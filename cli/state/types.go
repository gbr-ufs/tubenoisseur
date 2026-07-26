package state

// ProcessedVideo is a smaller version of [scraping.VideoEntry] for state.
// It only keeps the [scraping.VideoEntry.ID] for lookup, and [scraping.VideoEntry.Published]
// and [scraping.VideoEntry.Title] for generating the template text.
type ProcessedVideo struct {
	ID        string   `json:"id"`
	Published string   `json:"published"`
	Title     string   `json:"title"`
	URLs      []string `json:"urls"`
}

// ChannelState implements the JSON structure that represents the application
// state of a YouTube channel. It servers as the program's (external) memory.
type ChannelState struct {
	// ChannelHandle is the channel's identifier of its URL.
	ChannelHandle string `json:"channelHandle"`
	// DomainCounts consists of a domain-times-cited map.
	DomainCounts map[string]uint `json:"domainCounts"`
	// ExcludedDomains is a list of domains that shouldn't be registered
	// as a video's source.
	ExcludedDomains []string `json:"excludedDomains"`
	// FeedURL is the URL to the channel's RSS feed.
	FeedURL string `json:"feedURL"`
	// LastVideoID is the ID of the last video that was processed by the
	// program.
	LastVideoID string `json:"lastVideoID"`
	// Videos is the list of videos that were processed.
	Videos []ProcessedVideo `json:"videos"`
}
