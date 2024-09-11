package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/go-github/v62/github"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

// ghToken is the name of an environment variable whose value needs should be set with a GitHub personal access token
// (PAT). This PAT needs to have (at least) the 'repo' and 'notifications' [scopes]. If the notifications are inside an
// org that uses SAML SSO, the PAT must also be [authorised] for the org.
//
// [scopes]: https://docs.github.com/apps/building-oauth-apps/scopes-for-oauth-apps/
// [authorised]: https://docs.github.com/en/enterprise-cloud@latest/authentication/authenticating-with-saml-single-sign-on/authorizing-a-personal-access-token-for-use-with-saml-single-sign-on
const ghToken = "GITHUB_TOKEN"

const cmdName = "ginsu"

var requiredScopes = []string{"repo", "notifications"}

const (
	listPerPage     = 50
	loginDependabot = "dependabot[bot]"
)

// HTTP header keys.
const (
	// headerKeyScopes will list the scopes the token has authorised.
	headerKeyScopes = "X-Oauth-Scopes"

	// headerKeySAML will be populated if the token used is not authorised for [SAML SSO].
	//
	// [SAML SSO]: https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api?apiVersion=2022-11-28#personal-access-tokens-and-saml-sso
	headerKeySAML = "X-GitHub-SSO"
)

// Run these checks only once each.
var checkScopes, checkSAML sync.Once

// Exit statuses.
const (
	exitSuccess = iota
	exitUnknown
	exitHelp
	exitNoTokenSet
	exitTokenMissingScopes
	exitTokenMissingSAMLSSOAuth
)

// Flags.
var (
	flagDebug = pflag.BoolP("debug", "d", false,
		"show debugging output")
	flagHelp = pflag.BoolP("help", "h", false,
		"show help, and version from build info if available")
	flagOwnerAllowlist = pflag.StringSliceP("owner-allowlist", "o", []string{},
		"only drill down on these repo owners; comma-separated, not used if left unset")
)

// Static errors.
var (
	errOwnerNotOnAllowlist        = errors.New("repo owner not on allowlist")
	errTokenMissingRequiredScopes = errors.New("token does not have required scope")
	errTokenMissingSAMLSSOAuth    = errors.New("token is not authorised for SAML SSO with org(s)")
)

func main() {
	// Declare and defer this first, so that it runs last.
	// Assume there was an unknown error, unless we make it all the way to the bottom.
	exitStatus := exitUnknown
	defer func() {
		// Calling os.Exit needs to take place inside a closure so that the 'exitStatus' holding variable can be properly
		// accessed.
		os.Exit(exitStatus)
	}()

	pflag.Parse()

	if *flagHelp {
		helpOutput := &strings.Builder{}
		helpOutput.WriteString(cmdName)

		if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
			fmt.Fprintf(helpOutput, ", version %s\n\n", bi.Main.Version)
		}

		helpOutput.WriteString(pflag.CommandLine.FlagUsages())

		fmt.Print(helpOutput)

		exitStatus = exitHelp
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logOpts := log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: true,
	}

	if *flagDebug {
		logOpts.Level = log.DebugLevel
	} else {
		logOpts.Level = log.InfoLevel
	}

	// Declaring this here for re-use and potential future refactoring in case the log output needs to take a step back.
	var logOutput io.Writer = os.Stderr

	// For non-interactive runs, use a JSON logger.
	logOutFile, logOutIsFile := logOutput.(*os.File)
	if logOutIsFile && term.IsTerminal(int(logOutFile.Fd())) {
		logOpts.Formatter = log.TextFormatter
	} else {
		logOpts.Formatter = log.JSONFormatter
	}

	handler := log.NewWithOptions(logOutput, logOpts)
	slog.SetDefault(slog.New(handler))

	token, tokenSet := os.LookupEnv(ghToken)
	if !tokenSet {
		slog.Error("no GitHub token set",
			slog.String("env_var_key", ghToken),
		)

		exitStatus = exitNoTokenSet
		return
	}

	if err := run(ctx, token); err != nil {
		slog.Error("run",
			slog.Any("err", err),
		)

		switch {
		case errors.Is(err, errTokenMissingRequiredScopes):
			exitStatus = exitTokenMissingScopes

		case errors.Is(err, errTokenMissingSAMLSSOAuth):
			exitStatus = exitTokenMissingSAMLSSOAuth

		default:
			exitStatus = exitUnknown
		}

		return
	}

	exitStatus = exitSuccess
}

func run(ctx context.Context, token string) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	hc := oauth2.NewClient(ctx, ts)
	hc.Timeout = 5 * time.Second
	client := github.NewClient(hc)

	// Start sifting through notifications.
	firstPage, lastPage, err := listPageOfNotifications(ctx, client, 1)
	if err != nil {
		return fmt.Errorf("listing first page of notifications: %w", err)
	}

	slog.Debug("notification list pages",
		slog.Int("last", lastPage),
	)

	p := pool.NewWithResults[[]*github.Notification]().WithContext(ctx)

	for page := 2; page <= lastPage; page++ {
		page := page

		p.Go(func(ctx context.Context) ([]*github.Notification, error) {
			result, _, err := listPageOfNotifications(ctx, client, page)
			return result, err
		})
	}

	pagesAfterFirst, err := p.Wait()
	if err != nil {
		return fmt.Errorf("paginating notifications: %w", err)
	}

	notifications := slices.Concat(firstPage, slices.Concat(pagesAfterFirst...))

	slog.Debug("notifications",
		slog.Int("count", len(notifications)),
	)

	q := pool.New().WithErrors()
	dependabotCounter := &atomic.Uint64{}
	var i int

	for i = 0; i < len(notifications); i++ {
		index := i

		q.Go(func() error {
			return process(ctx, client, notifications[index], dependabotCounter)
		})
	}

	if err := q.Wait(); err != nil {
		return fmt.Errorf("working through notification pool: %w", err)
	}

	slog.Info("count of notifications processed",
		slog.Int("total", i),
		slog.Uint64("dependabot", dependabotCounter.Load()),
	)

	return nil
}

func checkTokenScopes(headers http.Header) error {
	scopesHeader, scopesHeaderExists := headers[http.CanonicalHeaderKey(headerKeyScopes)]
	if !scopesHeaderExists {
		return fmt.Errorf("%w: %s", errTokenMissingRequiredScopes, strings.Join(requiredScopes, ","))
	}

	scopesFound := map[string]struct{}{}

	slog.Debug("scopes header",
		slog.String("key", headerKeyScopes),
		slog.String("slice", fmt.Sprintf("%#v", scopesHeader)),
	)

	for index := range scopesHeader {
		for _, foundScope := range strings.Split(scopesHeader[index], ",") {
			trimmedScope := strings.TrimSpace(foundScope)
			scopesFound[trimmedScope] = struct{}{}
		}
	}

	for index := range requiredScopes {
		requiredScope := requiredScopes[index]

		if _, foundRequiredScope := scopesFound[requiredScope]; !foundRequiredScope {
			return fmt.Errorf("%w: %s", errTokenMissingRequiredScopes, requiredScope)
		}
	}

	return nil
}

func checkTokenSAML(headers http.Header) error {
	samlHeader, samlHeaderExists := headers[http.CanonicalHeaderKey(headerKeySAML)]
	if !samlHeaderExists {
		return nil
	}

	orgIDs := []string{}

	slog.Debug("SAML SSO header",
		slog.String("key", headerKeySAML),
		slog.String("slice", fmt.Sprintf("%#v", samlHeader)),
	)

	for i := range samlHeader {
		xValue := strings.Split(samlHeader[i], "; ")

		slog.Debug("header value",
			slog.Int("index", i),
			slog.String("slice", fmt.Sprintf("%#v", xValue)),
		)

		for j := range xValue {
			if !strings.HasPrefix(xValue[j], "organizations=") {
				continue
			}

			rawOrgIDs := strings.SplitAfter(xValue[j], "organizations=")

			slog.Debug("raw org IDs",
				slog.Int("index", j),
				slog.String("slice", fmt.Sprintf("%#v", rawOrgIDs)),
			)

			if len(rawOrgIDs) > 1 {
				orgIDs = append(orgIDs, strings.Split(rawOrgIDs[1], ",")...)
			}
		}
	}

	orgAPIURLs := []string{}

	for i := range orgIDs {
		fmtOrgAPIURL := "https://api.github.com/orgs/%s"
		orgAPIURLs = append(orgAPIURLs, fmt.Sprintf(fmtOrgAPIURL, orgIDs[i]))
	}

	return fmt.Errorf("%w: %s", errTokenMissingSAMLSSOAuth, strings.Join(orgAPIURLs, ";"))
}

func listPageOfNotifications(ctx context.Context, client *github.Client, page int) (
	[]*github.Notification, int, error,
) {
	opts := &github.NotificationListOptions{
		// If true, show notifications marked as read.
		All: false,

		// If true, only shows notifications in which the user is directly participating or mentioned.
		Participating: false,

		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: listPerPage,
		},
	}

	slog.Debug("started listing page of notifications",
		slog.Int("page_number", opts.Page),
	)

	defer slog.Debug("finished listing page of notifications",
		slog.Int("page_number", opts.Page),
	)

	nots, resp, err := client.Activity.ListNotifications(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("listing notifications: %w", err)
	}
	defer resp.Body.Close()

	slog.Debug("got page of notifications",
		slog.Int("page_number", opts.Page),
		slog.Int("count", len(nots)),
	)

	// Check the token we're using has the necessary scopes and SAML SSO auth, but only once each.
	var errTokenScopes, errSAML error

	checkScopes.Do(func() {
		errTokenScopes = checkTokenScopes(resp.Header)
	})

	if errTokenScopes != nil {
		return nil, 0, fmt.Errorf("checking token scopes: %w", errTokenScopes)
	}

	checkSAML.Do(func() {
		errSAML = checkTokenSAML(resp.Header)
	})

	if errSAML != nil {
		return nil, 0, errSAML
	}

	return nots, resp.LastPage, nil
}

func process(ctx context.Context, client *github.Client, ghn *github.Notification, dependabot *atomic.Uint64) error {
	slog.Debug("starting to process notification",
		slog.String("type", ghn.GetSubject().GetType()),
		slog.String("title", ghn.GetSubject().GetTitle()),
	)

	defer slog.Debug("finished processing notification",
		slog.String("type", ghn.GetSubject().GetType()),
		slog.String("title", ghn.GetSubject().GetTitle()),
	)

	switch ghn.GetSubject().GetType() {
	case "Issue":
		if err := lookAtIssue(ctx, client, ghn); err != nil {
			if !errors.Is(err, errOwnerNotOnAllowlist) {
				return err
			}

			slog.Warn("issue owner not on allowlist",
				slog.String("repo", ghn.GetRepository().GetFullName()),
				slog.String("title", ghn.GetSubject().GetTitle()),
				slog.String("type", ghn.GetSubject().GetType()),
				slog.String("url", ghn.GetSubject().GetURL()),
				slog.Time("updated_at", ghn.GetUpdatedAt().Time),
				slog.String("allowlist", fmt.Sprintf("%+v", *flagOwnerAllowlist)),
				slog.Any("err", err),
			)
		}

		return nil

	case "PullRequest":
		if err := lookAtPullRequest(ctx, client, ghn, dependabot); err != nil {
			if !errors.Is(err, errOwnerNotOnAllowlist) {
				return err
			}

			slog.Warn("PR owner not on allowlist",
				slog.String("repo", ghn.GetRepository().GetFullName()),
				slog.String("title", ghn.GetSubject().GetTitle()),
				slog.String("type", ghn.GetSubject().GetType()),
				slog.String("url", ghn.GetSubject().GetURL()),
				slog.Time("updated_at", ghn.GetUpdatedAt().Time),
				slog.String("allowlist", fmt.Sprintf("%+v", *flagOwnerAllowlist)),
				slog.Any("err", err),
			)
		}

		return nil

	default:
		slog.Warn("not an issue or a PR",
			slog.String("repo", ghn.GetRepository().GetFullName()),
			slog.String("type", ghn.GetSubject().GetType()),
			slog.String("title", ghn.GetSubject().GetTitle()),
		)

		return nil
	}
}

type details struct {
	owner, repo string
	number      int
}

func parseForDetails(ghn *github.Notification) (details, error) {
	xRepoFN := strings.Split(ghn.GetRepository().GetFullName(), "/")
	if len(xRepoFN) < 2 {
		return details{}, fmt.Errorf("repo full name '%s' did not split into at least two substrings",
			ghn.GetRepository().GetFullName())
	}

	owner := xRepoFN[0]
	repo := xRepoFN[1]

	if len(*flagOwnerAllowlist) > 0 && !slices.Contains(*flagOwnerAllowlist, owner) {
		return details{}, fmt.Errorf("%w: %s", errOwnerNotOnAllowlist, owner)
	}

	xURL := strings.Split(ghn.GetSubject().GetURL(), "/")
	if len(xURL) < 1 {
		return details{}, fmt.Errorf("notification subject URL '%s' did not split", ghn.GetSubject().GetURL())
	}

	subjectURLNumber := xURL[len(xURL)-1]

	number, err := strconv.Atoi(subjectURLNumber)
	if err != nil {
		return details{}, fmt.Errorf("could not convert subject URL number '%s': %w", subjectURLNumber, err)
	}

	return details{
		owner:  owner,
		repo:   repo,
		number: number,
	}, nil
}

func lookAtIssue(ctx context.Context, client *github.Client, ghn *github.Notification) error {
	slog.Debug("issue notification",
		slog.String("repo", ghn.GetRepository().GetFullName()),
		slog.String("title", ghn.GetSubject().GetTitle()),
		slog.String("type", ghn.GetSubject().GetType()),
		slog.String("url", ghn.GetSubject().GetURL()),
		slog.Time("updated_at", ghn.GetUpdatedAt().Time),
	)

	issDets, err := parseForDetails(ghn)
	if err != nil {
		return fmt.Errorf("parsing for details: %w", err)
	}

	issue, resp, err := client.Issues.Get(ctx, issDets.owner, issDets.repo, issDets.number)
	if err != nil {
		return fmt.Errorf("getting pull request: %w", err)
	}
	defer resp.Body.Close()

	slog.Info("state of issue",
		slog.Int("number", issue.GetNumber()),
		slog.String("state", issue.GetState()),
	)

	return nil
}

func lookAtPullRequest(ctx context.Context, client *github.Client, ghn *github.Notification, dependabot *atomic.Uint64,
) error {
	slog.Debug("PR notification",
		slog.String("repo", ghn.GetRepository().GetFullName()),
		slog.String("title", ghn.GetSubject().GetTitle()),
		slog.String("type", ghn.GetSubject().GetType()),
		slog.String("url", ghn.GetSubject().GetURL()),
		slog.Time("updated_at", ghn.GetUpdatedAt().Time),
	)

	prDets, err := parseForDetails(ghn)
	if err != nil {
		return fmt.Errorf("parsing for details: %w", err)
	}

	pr, resp, err := client.PullRequests.Get(ctx, prDets.owner, prDets.repo, prDets.number)
	if err != nil {
		return fmt.Errorf("getting pull request: %w", err)
	}
	defer resp.Body.Close()

	if pr.GetUser().GetLogin() == loginDependabot {
		slog.Debug("dependabot PR",
			slog.String("repo", pr.GetBase().GetRepo().GetFullName()),
			slog.String("title", pr.GetTitle()),
		)

		dependabot.Add(1)
	}

	if pr.GetState() != "closed" {
		slog.Debug("PR not closed, so leaving associated notification alone",
			slog.String("title", pr.GetTitle()),
			slog.String("user_login", pr.GetUser().GetLogin()),
			slog.String("repo", pr.GetBase().GetRepo().GetFullName()),
			slog.Int("number", prDets.number),
			slog.String("state", pr.GetState()),
		)

		return nil
	}

	slog.Info("PR is closed, marking as done",
		slog.String("title", pr.GetTitle()),
		slog.String("user_login", pr.GetUser().GetLogin()),
		slog.String("repo", pr.GetBase().GetRepo().GetFullName()),
		slog.Int("number", prDets.number),
		slog.String("state", pr.GetState()),
	)

	return markAsDone(ctx, client, ghn)
}

func markAsDone(ctx context.Context, client *github.Client, ghn *github.Notification) error {
	reqURL := client.BaseURL.String() + path.Join("notifications", "threads", ghn.GetID())

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("response status when attempting to mark as done: %s", resp.Status)
	}

	return nil
}
