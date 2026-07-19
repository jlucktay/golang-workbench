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
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Mode enum.Member[string]

var (
	ModeSolo = Mode{"solo"}
	ModeDuo  = Mode{"duo"}
	ModeTrio = Mode{"trio"}
	Modes    = enum.New(ModeSolo, ModeDuo, ModeTrio)
)

type Map enum.Member[string]

var (
	MapPerimeter      = Map{"perimeter"}
	MapDireMarsh      = Map{"dire marsh"}
	MapDireMarshNight = Map{"dire marsh (night)"}
	MapOutput         = Map{"output"}
	MapCryoArchive    = Map{"cryo archive"}
	Maps              = enum.New(MapPerimeter, MapDireMarsh, MapDireMarshNight, MapOutput, MapCryoArchive)
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
	if err := run(os.Stdout, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(out io.Writer, in io.Reader) error {
	ctx := context.Background()

	youtubeService, err := youtube.NewService(ctx, option.WithScopes(youtube.YoutubeScope))
	if err != nil {
		return fmt.Errorf("failed to create YouTube service: %w", err)
	}

	channelsCall := youtubeService.Channels.List([]string{"contentDetails"}).Mine(true)
	channelsResponse, err := channelsCall.Do()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	if len(channelsResponse.Items) == 0 {
		return fmt.Errorf("no channels found for the user")
	}

	uploadsPlaylistID := channelsResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads
	playlistItemsCall := youtubeService.PlaylistItems.List([]string{"snippet", "contentDetails"}).
		PlaylistId(uploadsPlaylistID).
		MaxResults(50)

	var allVideos []*youtube.Video
	pageToken := ""

	fmt.Fprintln(out, "Looking up oldest videos in Drafts collection...")

	for {
		if pageToken != "" {
			playlistItemsCall.PageToken(pageToken)
		}
		playlistResponse, err := playlistItemsCall.Do()
		if err != nil {
			return fmt.Errorf("failed to get playlist items: %w", err)
		}

		var videoIDs []string
		for _, item := range playlistResponse.Items {
			videoIDs = append(videoIDs, item.ContentDetails.VideoId)
		}

		if len(videoIDs) > 0 {
			videosCall := youtubeService.Videos.List([]string{"snippet", "status", "fileDetails", "recordingDetails"}).Id(strings.Join(videoIDs, ","))
			videosResponse, err := videosCall.Do()
			if err != nil {
				return fmt.Errorf("failed to get video details: %w", err)
			}

			for _, video := range videosResponse.Items {
				// Drafts usually have "private" status and maybe not published fully.
				// We identify "Drafts" as private videos for this workflow.
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

	var draftsToProcess []*youtube.Video
	if len(allVideos) > 10 {
		draftsToProcess = allVideos[:10]
	} else {
		draftsToProcess = allVideos
	}

	if len(draftsToProcess) == 0 {
		fmt.Fprintln(out, "No draft videos found.")
		return nil
	}

	var marathonPlaylistID string
	playlistsCall := youtubeService.Playlists.List([]string{"snippet"}).Mine(true).MaxResults(50)
	playlistsResp, err := playlistsCall.Do()
	if err == nil {
		for _, pl := range playlistsResp.Items {
			if pl.Snippet.Title == "My Marathon runs" {
				marathonPlaylistID = pl.Id
				break
			}
		}
	}

	if marathonPlaylistID == "" {
		fmt.Fprintln(out, "Creating 'My Marathon runs' playlist...")
		newPl := &youtube.Playlist{
			Snippet: &youtube.PlaylistSnippet{
				Title: "My Marathon runs",
			},
			Status: &youtube.PlaylistStatus{
				PrivacyStatus: "private",
			},
		}
		plResp, err := youtubeService.Playlists.Insert([]string{"snippet", "status"}, newPl).Do()
		if err != nil {
			fmt.Fprintf(out, "Warning: failed to create playlist: %v\n", err)
		} else {
			marathonPlaylistID = plResp.Id
		}
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

		var selectedMode string
		var selectedMap string
		var selectedShell string

		modeOptions := []huh.Option[string]{}
		for _, m := range Modes.Members() {
			modeOptions = append(modeOptions, huh.NewOption(m.Value, m.Value))
		}

		mapOptions := []huh.Option[string]{}
		for _, m := range Maps.Members() {
			mapOptions = append(mapOptions, huh.NewOption(m.Value, m.Value))
		}

		shellOptions := []huh.Option[string]{}
		for _, s := range Shells.Members() {
			shellOptions = append(shellOptions, huh.NewOption(s.Value, s.Value))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Mode").
					Options(modeOptions...).
					Value(&selectedMode),
				huh.NewSelect[string]().
					Title("Select Map").
					Options(mapOptions...).
					Value(&selectedMap),
				huh.NewSelect[string]().
					Title("Select Shell").
					Options(shellOptions...).
					Value(&selectedShell),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		tags := append(video.Snippet.Tags,
			fmt.Sprintf("video-games/marathon/mode/%s", selectedMode),
			fmt.Sprintf("video-games/marathon/map/%s", selectedMap),
			fmt.Sprintf("video-games/marathon/shell/%s", selectedShell),
		)
		video.Snippet.Tags = tags

		fileName := video.Snippet.Title
		if video.FileDetails != nil && video.FileDetails.FileName != "" {
			fileName = video.FileDetails.FileName
		}

		baseName := strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
		var recDate time.Time
		recDate, err = time.Parse("2006-01-02_15-04-05", baseName)
		if err == nil {
			if video.RecordingDetails == nil {
				video.RecordingDetails = &youtube.VideoRecordingDetails{}
			}
			video.RecordingDetails.RecordingDate = recDate.Format(time.RFC3339)
		} else {
			recDate, _ = time.Parse(time.RFC3339, video.Snippet.PublishedAt)
		}

		// YouTube Gaming category is 20
		video.Snippet.CategoryId = "20"

		// No API field for "AI was not used", but we ensure MadeForKids is false and no age restriction.
		if video.Status == nil {
			video.Status = &youtube.VideoStatus{}
		}
		video.Status.MadeForKids = false
		video.Status.SelfDeclaredMadeForKids = false
		video.Status.PrivacyStatus = "private"

		// Update title
		titleMode := strings.Title(selectedMode)
		titleShell := strings.Title(selectedShell)
		titleMap := strings.Title(selectedMap)
		video.Snippet.Title = fmt.Sprintf("%s %s on %s — %s", titleMode, titleShell, titleMap, recDate.Format("2006-01-02T15:04:05"))

		updateCall := youtubeService.Videos.Update([]string{"snippet", "status", "recordingDetails"}, video)
		_, err = updateCall.Do()
		if err != nil {
			fmt.Fprintf(out, "Failed to update metadata for video %s: %v\n", video.Id, err)
		} else {
			fmt.Fprintf(out, "Successfully updated metadata for video: %s\n", video.Id)
		}

		if marathonPlaylistID != "" {
			playlistItem := &youtube.PlaylistItem{
				Snippet: &youtube.PlaylistItemSnippet{
					PlaylistId: marathonPlaylistID,
					ResourceId: &youtube.ResourceId{
						Kind:    "youtube#video",
						VideoId: video.Id,
					},
				},
			}
			_, err = youtubeService.PlaylistItems.Insert([]string{"snippet"}, playlistItem).Do()
			if err != nil {
				fmt.Fprintf(out, "Note: couldn't add to playlist (might already be there or quota error): %v\n", err)
			} else {
				fmt.Fprintf(out, "Added to 'My Marathon runs' playlist.\n")
			}
		}
	}

	fmt.Fprintln(out, "All done!")
	return nil
}
