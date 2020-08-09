package main

type VideosResult struct {
	Error                string    `json:"error"`
	Limit                int       `json:"limit"`
	Offset               int       `json:"offset"`
	NumberOfPageResults  int       `json:"number_of_page_results"`
	NumberOfTotalResults int       `json:"number_of_total_results"`
	StatusCode           int       `json:"status_code"`
	Results              []Results `json:"results"`
}

type Results struct {
	APIDetailURL    string            `json:"api_detail_url"`
	Associations    []Associations    `json:"associations"`
	Deck            string            `json:"deck"`
	EmbedPlayer     string            `json:"embed_player"`
	GUID            string            `json:"guid"`
	ID              int               `json:"id"`
	LengthSeconds   int               `json:"length_seconds"`
	Name            string            `json:"name"`
	Premium         bool              `json:"premium"`
	PublishDate     string            `json:"publish_date"`
	SiteDetailURL   string            `json:"site_detail_url"`
	Image           Image             `json:"image"`
	User            string            `json:"user"`
	Hosts           string            `json:"hosts"`
	Crew            string            `json:"crew"`
	VideoType       string            `json:"video_type"`
	VideoShow       VideoShow         `json:"video_show"`
	VideoCategories []VideoCategories `json:"video_categories"`
	SavedTime       interface{}       `json:"saved_time"`
	YoutubeID       string            `json:"youtube_id"`
	LowURL          string            `json:"low_url"`
	HighURL         string            `json:"high_url"`
	HdURL           string            `json:"hd_url"`
	URL             string            `json:"url"`
}

type Associations struct {
	APIDetailURL  string `json:"api_detail_url"`
	GUID          string `json:"guid"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SiteDetailURL string `json:"site_detail_url"`
}

type Image struct {
	IconURL        string `json:"icon_url"`
	MediumURL      string `json:"medium_url"`
	ScreenURL      string `json:"screen_url"`
	ScreenLargeURL string `json:"screen_large_url"`
	SmallURL       string `json:"small_url"`
	SuperURL       string `json:"super_url"`
	ThumbURL       string `json:"thumb_url"`
	TinyURL        string `json:"tiny_url"`
	OriginalURL    string `json:"original_url"`
	ImageTags      string `json:"image_tags"`
}

type VideoShow struct {
	APIDetailURL  string `json:"api_detail_url"`
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Position      int    `json:"position"`
	SiteDetailURL string `json:"site_detail_url"`
	Image         Image  `json:"image"`
	Logo          Logo   `json:"logo"`
}

type Logo struct {
	IconURL        string `json:"icon_url"`
	MediumURL      string `json:"medium_url"`
	ScreenURL      string `json:"screen_url"`
	ScreenLargeURL string `json:"screen_large_url"`
	SmallURL       string `json:"small_url"`
	SuperURL       string `json:"super_url"`
	ThumbURL       string `json:"thumb_url"`
	TinyURL        string `json:"tiny_url"`
	OriginalURL    string `json:"original_url"`
}

type VideoCategories struct {
	APIDetailURL  string `json:"api_detail_url"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SiteDetailURL string `json:"site_detail_url"`
}
