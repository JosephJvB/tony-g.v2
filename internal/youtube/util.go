package youtube

import (
	"net/url"
	"strings"
)

func GetYoutubeVideoID(urlStr string) string {
	yturl, err := url.Parse(strings.TrimSpace(urlStr))
	if err != nil {
		return ""
	}

	// handles:
	// https://www.youtube
	// https://music.youtube
	// https://youtube
	if strings.HasSuffix(yturl.Host, "youtube.com") {
		return yturl.Query().Get("v")
	}

	if strings.HasSuffix(yturl.Host, "youtu.be") {
		return yturl.Path[1:] // rm leading slash
	}

	return ""
}

func GetReviewVideos(items []PlaylistItem) []PlaylistItem {
	filtered := []PlaylistItem{}

	for _, item := range items {
		if item.Status.PrivacyStatus == "private" {
			continue
		}
		if item.Snippet.VideoOwnerChannelId != item.Snippet.ChannelId {
			continue
		}
		// album review, ep review, compilation review, mixtape review
		// why are these videos in his playlist? chaotic
		if strings.HasSuffix(strings.TrimSpace(item.Snippet.Title), "REVIEW") {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}
