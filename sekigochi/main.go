package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/orsinium-labs/enum"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Mode enum.Member[string]

var (
	ModePvPvE = Mode{"pvpve"}
	ModePvE   = Mode{"pve"}
	Modes     = enum.New(ModePvPvE, ModePvE)
)

type Crew enum.Member[string]

var (
	CrewSolo = Crew{"solo"}
	CrewDuo  = Crew{"duo"}
	CrewTrio = Crew{"trio"}
	Crews    = enum.New(CrewSolo, CrewDuo, CrewTrio)
)

type Map enum.Member[string]

var (
	MapPerimeter      = Map{"perimeter"}
	MapDireMarsh      = Map{"dire marsh"}
	MapDireMarshNight = Map{"dire marsh (night)"}
	MapOutpost        = Map{"outpost"}
	MapCryoArchive    = Map{"cryo archive"}
	Maps              = enum.New(MapPerimeter, MapDireMarsh, MapDireMarshNight, MapOutpost, MapCryoArchive)
)

type Shell enum.Member[string]

var (
	ShellDestroyer = Shell{"destroyer"}
	ShellTriage    = Shell{"triage"}
	ShellRecon     = Shell{"recon"}
	ShellThief     = Shell{"thief"}
	ShellAssassin  = Shell{"assassin"}
	ShellVandal    = Shell{"vandal"}
	ShellSentinel  = Shell{"sentinel"}
	Shells         = enum.New(ShellDestroyer, ShellTriage, ShellRecon, ShellThief, ShellAssassin, ShellVandal, ShellSentinel)
)

func main() {
	app := &cli.App{
		Name:  "sekigochi",
		Usage: "YouTube Marathon video updater",
		Action: func(c *cli.Context) error {
			return run(c.Context, os.Stdout, os.Stdin)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, out io.Writer, in io.Reader) error {
	youtubeService, err := youtube.NewService(ctx, option.WithScopes(youtube.YoutubeScope))
	if err != nil {
		return fmt.Errorf("failed to create YouTube service: %w", err)
	}

	fmt.Fprintln(out, "Looking up oldest videos in Drafts collection...")
	draftsToProcess, err := fetchOldestDrafts(youtubeService, 10)
	if err != nil {
		return fmt.Errorf("failed to fetch drafts: %w", err)
	}

	if len(draftsToProcess) == 0 {
		fmt.Fprintln(out, "No draft videos found.")
		return nil
	}

	duplicates := findDuplicates(draftsToProcess)
	if len(duplicates) > 0 {
		fmt.Fprintf(out, "WARNING: Found duplicate video names: %v\n", duplicates)
	}

	marathonPlaylistID, err := getOrCreatePlaylist(youtubeService, "My Marathon runs", out)
	if err != nil {
		fmt.Fprintf(out, "Warning: couldn't ensure playlist: %v\n", err)
	}

	for i, video := range draftsToProcess {
		fmt.Fprintf(out, "\n[%d/%d] Processing video: %s (ID: %s)\n", i+1, len(draftsToProcess), video.Snippet.Title, video.Id)

		var openInFirefox bool
		err := huh.NewConfirm().
			Title("Open this video in Firefox?").
			Value(&openInFirefox).
			Run()
		if err != nil {
			return err
		}

		if openInFirefox {
			vidURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id)
			exec.Command("firefox", vidURL).Start()
		}

		var selectedMode, selectedCrew, selectedMap, selectedShell string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Mode").
					Options(mapOptions(Modes.Members(), func(m Mode) string { return m.Value })...).
					Value(&selectedMode),
				huh.NewSelect[string]().
					Title("Select Crew").
					Options(mapOptions(Crews.Members(), func(c Crew) string { return c.Value })...).
					Value(&selectedCrew),
				huh.NewSelect[string]().
					Title("Select Map").
					Options(mapOptions(Maps.Members(), func(m Map) string { return m.Value })...).
					Value(&selectedMap),
				huh.NewSelect[string]().
					Title("Select Shell").
					Options(mapOptions(Shells.Members(), func(s Shell) string { return s.Value })...).
					Value(&selectedShell),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		fileName := video.Snippet.Title
		if video.FileDetails != nil && video.FileDetails.FileName != "" {
			fileName = video.FileDetails.FileName
		}
		recDate := parseRecordingDate(fileName, video.Snippet.PublishedAt)

		updateVideoMetadata(video, selectedCrew, selectedShell, selectedMap, selectedMode, recDate)

		updateCall := youtubeService.Videos.Update([]string{"snippet", "status", "recordingDetails"}, video)
		if _, err = updateCall.Do(); err != nil {
			fmt.Fprintf(out, "Failed to update metadata for video %s: %v\n", video.Id, err)
		} else {
			fmt.Fprintf(out, "Successfully updated metadata for video: %s\n", video.Id)
		}

		if marathonPlaylistID != "" {
			err = addToPlaylist(youtubeService, marathonPlaylistID, video.Id)
			if err != nil {
				fmt.Fprintf(out, "Note: couldn't add to playlist (might already be there): %v\n", err)
			} else {
				fmt.Fprintf(out, "Added to 'My Marathon runs' playlist.\n")
			}
		}
	}

	fmt.Fprintln(out, "All done!")
	return nil
}

// -----------------------------------------------------------------------------
// YouTube API Helpers
// -----------------------------------------------------------------------------

func fetchOldestDrafts(service *youtube.Service, limit int) ([]*youtube.Video, error) {
	channelsCall := service.Channels.List([]string{"contentDetails"}).Mine(true)
	channelsResponse, err := channelsCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	if len(channelsResponse.Items) == 0 {
		return nil, fmt.Errorf("no channels found for the user")
	}

	uploadsPlaylistID := channelsResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads
	playlistItemsCall := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
		PlaylistId(uploadsPlaylistID).
		MaxResults(50)

	var allVideos []*youtube.Video
	pageToken := ""

	for {
		if pageToken != "" {
			playlistItemsCall.PageToken(pageToken)
		}
		playlistResponse, err := playlistItemsCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist items: %w", err)
		}

		var videoIDs []string
		for _, item := range playlistResponse.Items {
			videoIDs = append(videoIDs, item.ContentDetails.VideoId)
		}

		if len(videoIDs) > 0 {
			videosCall := service.Videos.List([]string{"snippet", "status", "fileDetails", "recordingDetails"}).
				Id(strings.Join(videoIDs, ","))
			videosResponse, err := videosCall.Do()
			if err != nil {
				return nil, fmt.Errorf("failed to get video details: %w", err)
			}

			for _, video := range videosResponse.Items {
				if video.Status != nil && video.Status.PrivacyStatus == "private" {
					allVideos = append(allVideos, video)
				}
			}
		}

		pageToken = playlistResponse.NextPageToken
		if pageToken == "" {
			break
		}
	}

	sort.Slice(allVideos, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, allVideos[i].Snippet.PublishedAt)
		t2, err2 := time.Parse(time.RFC3339, allVideos[j].Snippet.PublishedAt)
		if err1 != nil || err2 != nil {
			return false
		}
		return t1.Before(t2)
	})

	if len(allVideos) > limit {
		return allVideos[:limit], nil
	}
	return allVideos, nil
}

func getOrCreatePlaylist(service *youtube.Service, title string, out io.Writer) (string, error) {
	playlistsCall := service.Playlists.List([]string{"snippet"}).Mine(true).MaxResults(50)
	playlistsResp, err := playlistsCall.Do()
	if err == nil {
		for _, pl := range playlistsResp.Items {
			if pl.Snippet.Title == title {
				return pl.Id, nil
			}
		}
	}

	fmt.Fprintf(out, "Creating '%s' playlist...\n", title)
	newPl := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{Title: title},
		Status:  &youtube.PlaylistStatus{PrivacyStatus: "private"},
	}
	plResp, err := service.Playlists.Insert([]string{"snippet", "status"}, newPl).Do()
	if err != nil {
		return "", err
	}
	return plResp.Id, nil
}

func addToPlaylist(service *youtube.Service, playlistID, videoID string) error {
	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}
	_, err := service.PlaylistItems.Insert([]string{"snippet"}, playlistItem).Do()
	return err
}

func updateVideoMetadata(video *youtube.Video, crew, shell, mapName, mode string, date time.Time) {
	video.Snippet.Tags = append(video.Snippet.Tags,
		fmt.Sprintf("video-games/marathon/mode/%s", mode),
		fmt.Sprintf("video-games/marathon/crew/%s", crew),
		fmt.Sprintf("video-games/marathon/map/%s", mapName),
		fmt.Sprintf("video-games/marathon/shell/%s", shell),
	)

	video.Snippet.Title = formatTitle(crew, shell, mapName, mode, date)
	video.Snippet.CategoryId = "20"

	if video.RecordingDetails == nil {
		video.RecordingDetails = &youtube.VideoRecordingDetails{}
	}
	video.RecordingDetails.RecordingDate = date.Format(time.RFC3339)

	if video.Status == nil {
		video.Status = &youtube.VideoStatus{}
	}
	video.Status.MadeForKids = false
	video.Status.SelfDeclaredMadeForKids = false
	video.Status.PrivacyStatus = "private"
}

// -----------------------------------------------------------------------------
// UI & Logic Helpers
// -----------------------------------------------------------------------------

func mapOptions[T any](items []T, get func(T) string) []huh.Option[string] {
	opts := make([]huh.Option[string], len(items))
	for i, item := range items {
		val := get(item)
		opts[i] = huh.NewOption(val, val)
	}
	return opts
}

func findDuplicates(videos []*youtube.Video) []string {
	seen := make(map[string]bool)
	var duplicates []string
	for _, v := range videos {
		name := v.Snippet.Title
		if v.FileDetails != nil && v.FileDetails.FileName != "" {
			name = v.FileDetails.FileName
		}
		if seen[name] {
			duplicates = append(duplicates, name)
		}
		seen[name] = true
	}
	return duplicates
}

func parseRecordingDate(fileName, fallback string) time.Time {
	baseName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	recDate, err := time.Parse("2006-01-02_15-04-05", baseName)
	if err != nil {
		recDate, _ = time.Parse(time.RFC3339, fallback)
	}
	return recDate
}

func formatTitle(crew, shell, mapName, mode string, date time.Time) string {
	var titleMode string
	switch strings.ToLower(mode) {
	case "pve":
		titleMode = "PvE"
	case "pvpve":
		titleMode = "PvPvE"
	default:
		titleMode = strings.Title(mode)
	}
	return fmt.Sprintf("%s %s on %s (%s) — %s", strings.Title(crew), strings.Title(shell), strings.Title(mapName), titleMode, date.Format("2006-01-02T15:04:05"))
}
