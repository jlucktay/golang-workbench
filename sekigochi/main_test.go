package main

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"google.golang.org/api/youtube/v3"
)

func TestFindDuplicates(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		videos   []*youtube.Video
		expected []string
	}{
		{
			name: "no duplicates",
			videos: []*youtube.Video{
				{Snippet: &youtube.VideoSnippet{Title: "vid1"}},
				{Snippet: &youtube.VideoSnippet{Title: "vid2"}},
			},
			expected: nil,
		},
		{
			name: "duplicates by title",
			videos: []*youtube.Video{
				{Snippet: &youtube.VideoSnippet{Title: "vid1"}},
				{Snippet: &youtube.VideoSnippet{Title: "vid1"}},
			},
			expected: []string{"vid1"},
		},
		{
			name: "duplicates by file name",
			videos: []*youtube.Video{
				{FileDetails: &youtube.VideoFileDetails{FileName: "file1.mp4"}},
				{FileDetails: &youtube.VideoFileDetails{FileName: "file1.mp4"}},
				{FileDetails: &youtube.VideoFileDetails{FileName: "file2.mp4"}},
			},
			expected: []string{"file1.mp4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got := findDuplicates(tt.videos)
			is.Equal(got, tt.expected)
		})
	}
}

func TestParseRecordingDate(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		fileName string
		fallback string
		expected time.Time
	}{
		{
			name:     "valid filename date",
			fileName: "2026-07-12_07-13-40.mp4",
			fallback: "2025-01-01T00:00:00Z",
			expected: time.Date(2026, 7, 12, 7, 13, 40, 0, time.UTC),
		},
		{
			name:     "invalid filename date falls back",
			fileName: "some-video.mp4",
			fallback: "2025-01-01T00:00:00Z",
			expected: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got := parseRecordingDate(tt.fileName, tt.fallback)
			is.Equal(got, tt.expected)
		})
	}
}

func TestFormatTitle(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		crew     string
		shell    string
		mapName  string
		mode     string
		date     time.Time
		expected string
	}{
		{
			name:     "pvpve title",
			crew:     "solo",
			shell:    "sentinel",
			mapName:  "outpost",
			mode:     "pvpve",
			date:     time.Date(2026, 7, 12, 7, 13, 40, 0, time.UTC),
			expected: "Solo Sentinel on Outpost (PvPvE) — 2026-07-12T07:13:40",
		},
		{
			name:     "pve title",
			crew:     "duo",
			shell:    "triage",
			mapName:  "dire marsh (night)",
			mode:     "pve",
			date:     time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "Duo Triage on Dire Marsh (Night) (PvE) — 2025-01-01T12:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			got := formatTitle(tt.crew, tt.shell, tt.mapName, tt.mode, tt.date)
			is.Equal(got, tt.expected)
		})
	}
}
