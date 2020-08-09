package main

type VideosResult struct {
	Error                string    `json:"error,omitempty"`
	Limit                int       `json:"limit,omitempty"`
	Offset               int       `json:"offset,omitempty"`
	NumberOfPageResults  int       `json:"number_of_page_results,omitempty"`
	NumberOfTotalResults int       `json:"number_of_total_results,omitempty"`
	StatusCode           int       `json:"status_code,omitempty"`
	Results              []Results `json:"results,omitempty"`
}

type Results struct {
	APIDetailURL    string            `json:"api_detail_url,omitempty"`
	Associations    []Associations    `json:"associations,omitempty"`
	Deck            string            `json:"deck,omitempty"`
	EmbedPlayer     string            `json:"embed_player,omitempty"`
	GUID            string            `json:"guid,omitempty"`
	ID              int               `json:"id,omitempty"`
	LengthSeconds   int               `json:"length_seconds,omitempty"`
	Name            string            `json:"name,omitempty"`
	Premium         bool              `json:"premium,omitempty"`
	PublishDate     string            `json:"publish_date,omitempty"`
	SiteDetailURL   string            `json:"site_detail_url,omitempty"`
	Image           *Image            `json:"image,omitempty"`
	User            string            `json:"user,omitempty"`
	Hosts           string            `json:"hosts,omitempty"`
	Crew            string            `json:"crew,omitempty"`
	VideoType       string            `json:"video_type,omitempty"`
	VideoShow       *VideoShow        `json:"video_show,omitempty"`
	VideoCategories []VideoCategories `json:"video_categories,omitempty"`
	SavedTime       interface{}       `json:"saved_time,omitempty"`
	YoutubeID       string            `json:"youtube_id,omitempty"`
	LowURL          string            `json:"low_url,omitempty"`
	HighURL         string            `json:"high_url,omitempty"`
	HdURL           string            `json:"hd_url,omitempty"`
	URL             string            `json:"url,omitempty"`
}

type Associations struct {
	APIDetailURL  string `json:"api_detail_url,omitempty"`
	GUID          string `json:"guid,omitempty"`
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
}

type Image struct {
	IconURL        string `json:"icon_url,omitempty"`
	MediumURL      string `json:"medium_url,omitempty"`
	ScreenURL      string `json:"screen_url,omitempty"`
	ScreenLargeURL string `json:"screen_large_url,omitempty"`
	SmallURL       string `json:"small_url,omitempty"`
	SuperURL       string `json:"super_url,omitempty"`
	ThumbURL       string `json:"thumb_url,omitempty"`
	TinyURL        string `json:"tiny_url,omitempty"`
	OriginalURL    string `json:"original_url,omitempty"`
	ImageTags      string `json:"image_tags,omitempty"`
}

type VideoShow struct {
	APIDetailURL  string `json:"api_detail_url,omitempty"`
	ID            int    `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Position      int    `json:"position,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
	Image         *Image `json:"image,omitempty"`
	Logo          *Logo  `json:"logo,omitempty"`
}

type Logo struct {
	IconURL        string `json:"icon_url,omitempty"`
	MediumURL      string `json:"medium_url,omitempty"`
	ScreenURL      string `json:"screen_url,omitempty"`
	ScreenLargeURL string `json:"screen_large_url,omitempty"`
	SmallURL       string `json:"small_url,omitempty"`
	SuperURL       string `json:"super_url,omitempty"`
	ThumbURL       string `json:"thumb_url,omitempty"`
	TinyURL        string `json:"tiny_url,omitempty"`
	OriginalURL    string `json:"original_url,omitempty"`
}

type VideoCategories struct {
	APIDetailURL  string `json:"api_detail_url,omitempty"`
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	SiteDetailURL string `json:"site_detail_url,omitempty"`
}
