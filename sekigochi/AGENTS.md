# `sekigochi` CLI Tool — Ground Rules

## Code structure

- The `main` entry point of the CLI tool will be a wrapper around a discrete and fully-testable `run` func.

## Libraries and dependencies

- Interactive CLI prompts will be built and run with the `github.com/charmbracelet` family of packages.
- YouTube API client directly from Google, no third-party go-betweens.
- Enums are implemented with `github.com/orsinium-labs/enum`.

## The flow of the CLI tool

  1. Look up the 10 oldest videos in my YouTube Drafts collection.
  1. Prompt to optionally open the oldest video to play in Firefox.
  1. Update the metadata on that same video as follows:
    1. Add the video to the `My Marathon runs` playlist.
    1. Set the `Audience` field to `No, it's not 'Made for Kids'`.
    1. Do not restrict the video to viewers over 18 only.
    1. AU was not used to generate or edit the content.
    1. Set one tag for each of the following on the video, in the format `video-games/marathon/<key>/<value>` (all lower-cased), with CLI prompts to choose values from a list of enums:
      1. Mode: Solo, Duo, Trio
      1. Map: Perimeter, Dire Marsh, Dire Marsh (Night), Output, Cryo Archive
      1. Shell: Destroyer, Triage, Recon, Thief, Assassin, Vandal, Sentinel
    1. Set the recording date, based on the video's original filename which is in the format `YYYY-MM-DD_HH-MM-SS.mp4`.
    1. Under category, set the game title to `Marathon (2026)`.
    1. Change the title to `<mode> <shell> on <map> — YYYY-MM-DDTHH:MM:SS`.
    1. Set visibility to `Private`.
  1. Repeat this Firefox open option and metadata update process for each video in the fetched list.
