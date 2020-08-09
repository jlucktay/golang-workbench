package main

type VideosResult struct {
	Limit                int `json:"limit,omitempty"`
	NumberOfPageResults  int `json:"number_of_page_results,omitempty"`
	NumberOfTotalResults int `json:"number_of_total_results,omitempty"`
	Offset               int `json:"offset,omitempty"`
	StatusCode           int `json:"status_code,omitempty"`

	Error string `json:"error,omitempty"`

	Results []Video `json:"results,omitempty"`
}

type Video struct {
	ID            int `json:"id,omitempty"`
	LengthSeconds int `json:"length_seconds,omitempty"`

	Premium bool `json:"premium,omitempty"`

	APIDetailURL  string `json:"api_detail_url,omitempty"`
	Crew          string `json:"crew,omitempty"`
	Deck          string `json:"deck,omitempty"`
	EmbedPlayer   string `json:"embed_player,omitempty"`
	GUID          string `json:"guid,omitempty"`
	HdURL         string `json:"hd_url,omitempty"`
	HighURL       string `json:"high_url,omitempty"`
	Hosts         string `json:"hosts,omitempty"`
	LowURL        string `json:"low_url,omitempty"`
	Name          string `json:"name,omitempty"`
	PublishDate   string `json:"publish_date,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
	URL           string `json:"url,omitempty"`
	User          string `json:"user,omitempty"`
	VideoType     string `json:"video_type,omitempty"`
	YoutubeID     string `json:"youtube_id,omitempty"`

	Associations    []Associations    `json:"associations,omitempty"`
	Image           *Image            `json:"image,omitempty"`
	SavedTime       interface{}       `json:"saved_time,omitempty"`
	VideoCategories []VideoCategories `json:"video_categories,omitempty"`
	VideoShow       *VideoShow        `json:"video_show,omitempty"`
}

type Associations struct {
	ID int `json:"id,omitempty"`

	APIDetailURL  string `json:"api_detail_url,omitempty"`
	GUID          string `json:"guid,omitempty"`
	Name          string `json:"name,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
}

type Image struct {
	IconURL        string `json:"icon_url,omitempty"`
	ImageTags      string `json:"image_tags,omitempty"`
	MediumURL      string `json:"medium_url,omitempty"`
	OriginalURL    string `json:"original_url,omitempty"`
	ScreenLargeURL string `json:"screen_large_url,omitempty"`
	ScreenURL      string `json:"screen_url,omitempty"`
	SmallURL       string `json:"small_url,omitempty"`
	SuperURL       string `json:"super_url,omitempty"`
	ThumbURL       string `json:"thumb_url,omitempty"`
	TinyURL        string `json:"tiny_url,omitempty"`
}

type VideoShow struct {
	ID       int `json:"id,omitempty"`
	Position int `json:"position,omitempty"`

	APIDetailURL  string `json:"api_detail_url,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
	Title         string `json:"title,omitempty"`

	Image *Image `json:"image,omitempty"`
	Logo  *Logo  `json:"logo,omitempty"`
}

type Logo struct {
	IconURL        string `json:"icon_url,omitempty"`
	MediumURL      string `json:"medium_url,omitempty"`
	OriginalURL    string `json:"original_url,omitempty"`
	ScreenLargeURL string `json:"screen_large_url,omitempty"`
	ScreenURL      string `json:"screen_url,omitempty"`
	SmallURL       string `json:"small_url,omitempty"`
	SuperURL       string `json:"super_url,omitempty"`
	ThumbURL       string `json:"thumb_url,omitempty"`
	TinyURL        string `json:"tiny_url,omitempty"`
}

type VideoCategories struct {
	ID int `json:"id,omitempty"`

	APIDetailURL  string `json:"api_detail_url,omitempty"`
	Name          string `json:"name,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
}
