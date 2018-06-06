package webhook

// Service ..
type Service struct {
	environment string
	hostname    string
}

// StanRequest describes the payload received in the webhook endpoint
type StanRequest struct {
	Payload      []*Payload `json:"payload"`
	Skip         int        `json:"skip"`
	Take         int        `json:"take"`
	TotalRecords int        `json:"totalRecords"`
}

// Payload describes entries in the Payload array received in the request
type Payload struct {
	Country      string `json:"country,omitempty"`
	Description  string `json:"description,omitempty"`
	DRM          bool   `json:"drm,omitempty"`
	EpisodeCount int    `json:"episodeCount,omitempty"`
	Genre        string `json:"genre,omitempty"`
	Image        struct {
		ShowImage string `json:"showImage,omitempty"`
	} `json:"image,omitempty"`
	Language      string      `json:"language,omitempty"`
	NextEpisode   interface{} `json:"nextEpisode,omitempty"`
	PrimaryColour string      `json:"primaryColour,omitempty"`
	Seasons       []struct {
		Slug string `json:"slug,omitempty"`
	} `json:"seasons,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Title     string `json:"title,omitempty"`
	TvChannel string `json:"tvChannel,omitempty"`
}

// NextEpisode describes info in the field
type NextEpisode struct {
	Channel     interface{} `json:"channel"`
	ChannelLogo string      `json:"channelLogo"`
	Date        interface{} `json:"date"`
	HTML        string      `json:"html"`
	URL         string      `json:"url"`
}

// StanResponse describes the expected response to be sent to Stan
type StanResponse struct {
	Response []*Response `json:"response,omitempty"`
}

// Response describes expected entry per array item in the StanResponse
type Response struct {
	Image string `json:"image,omitempty"`
	Slug  string `json:"slug,omitempty"`
	Title string `json:"title,omitempty"`
}
