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
