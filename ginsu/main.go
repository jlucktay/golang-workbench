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
	"slices"
	"strconv"
	"strings"
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

const listPerPage = 50

const (
	exitSuccess = iota
	exitUnknown
	exitNoTokenSet
)

var (
	flagDebug = pflag.BoolP("debug", "d", false,
		"show debugging output")
	flagOwnerAllowlist = pflag.StringSliceP("owner-allowlist", "o", []string{},
		"only drill down on these repo owners; comma-separated, not used if left unset")
)

var errOwnerNotOnAllowlist = errors.New("repo owner not on allowlist")

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

		exitStatus = exitUnknown
		return
	}

	exitStatus = exitSuccess
}

func run(ctx context.Context, token string) error {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	hc := oauth2.NewClient(ctx, ts)
	hc.Timeout = 5 * time.Second
	client := github.NewClient(hc)

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

	for i := 0; i < len(notifications); i++ {
		index := i

		q.Go(func() error {
			return process(ctx, client, notifications[index])
		})
	}

	return q.Wait()
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

	return nots, resp.LastPage, nil
}

func process(ctx context.Context, client *github.Client, ghn *github.Notification) error {
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

			slog.Warn("owner not on allowlist",
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
		return lookAtPullRequest(ctx, client, ghn)

	default:
		slog.Warn("not an issue or a PR",
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

func lookAtPullRequest(ctx context.Context, client *github.Client, ghn *github.Notification) error {
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

	if pr.GetState() != "closed" {
		slog.Debug("PR not closed, so leaving associated notification alone",
			slog.String("repo", ghn.GetRepository().GetFullName()),
			slog.Int("number", prDets.number),
			slog.String("title", pr.GetTitle()),
			slog.String("state", pr.GetState()),
		)

		return nil
	}

	slog.Info("PR is closed, marking as done",
		slog.String("repo", ghn.GetRepository().GetFullName()),
		slog.Int("number", prDets.number),
		slog.String("title", pr.GetTitle()),
		slog.String("state", pr.GetState()),
	)

	return markAsDone(ctx, client, ghn)
}

func markAsDone(ctx context.Context, client *github.Client, ghn *github.Notification) error {
	reqURL := "https://api.github.com/notifications/threads/" + ghn.GetID()

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
