package main

type FacebookExport struct {
	Title      string `json:"title"`
	ThreadType string `json:"thread_type"`
	ThreadPath string `json:"thread_path"`

	IsStillParticipant bool `json:"is_still_participant"`

	Messages     []FacebookMessage     `json:"messages"`
	Participants []FacebookParticipant `json:"participants"`
}

type FacebookParticipant struct {
	Name string `json:"name"`
}

type FacebookMessage struct {
	SenderName  string              `json:"sender_name"`
	TimestampMs int64               `json:"timestamp_ms"`
	Type        FacebookMessageType `json:"type"`

	Content string `json:"content,omitempty"`
	IP      string `json:"ip,omitempty"`

	CallDuration int `json:"call_duration,omitempty"`

	Missed bool `json:"missed,omitempty"`

	Share   FacebookMessageShare   `json:"share,omitempty"`
	Sticker FacebookMessageSticker `json:"sticker,omitempty"`

	Files     []FacebookMessageFile     `json:"files,omitempty"`
	Photos    []FacebookMessagePhoto    `json:"photos,omitempty"`
	Reactions []FacebookMessageReaction `json:"reactions,omitempty"`
	Videos    []FacebookMessageVideo    `json:"videos,omitempty"`
}

type FacebookMessageType string

const (
	FMTCall    FacebookMessageType = "Call"
	FMTGeneric FacebookMessageType = "Generic"
	FMTShare   FacebookMessageType = "Share"
)

type FacebookMessageSticker struct {
	URI string `json:"uri"`
}

type FacebookMessageShare struct {
	Link string `json:"link"`
}

type FacebookMessagePhoto struct {
	CreationTimestamp int    `json:"creation_timestamp"`
	URI               string `json:"uri"`
}

type FacebookMessageVideo struct {
	CreationTimestamp int                           `json:"creation_timestamp"`
	Thumbnail         FacebookMessageVideoThumbnail `json:"thumbnail"`
	URI               string                        `json:"uri"`
}

type FacebookMessageVideoThumbnail struct {
	URI string `json:"uri"`
}

type FacebookMessageReaction struct {
	Actor    string `json:"actor"`
	Reaction string `json:"reaction"`
}

type FacebookMessageFile struct {
	CreationTimestamp int    `json:"creation_timestamp"`
	URI               string `json:"uri"`
}
