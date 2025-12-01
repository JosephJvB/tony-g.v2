package youtube

import (
	"net/url"
	"regexp"
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

// from https://github.com/JosephJvB/tony-g/blob/main/internal/spotify/models.go
func CleanSongTitle(songTitle string) string {
	// Apple music main cases are (feat. ...) and [feat. ...]
	rmParens := regexp.MustCompile(`\\*\(feat.[^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[feat.[^)]*\]*`)
	songTitle = rmParens.ReplaceAllLiteralString(songTitle, "")
	songTitle = rmSquareBrackets.ReplaceAllLiteralString(songTitle, "")

	// Youtube description titles
	rmFtDot := regexp.MustCompile(`\\*( ft\..*)`)
	songTitle = rmFtDot.ReplaceAllLiteralString(songTitle, "")
	rmFeatDot := regexp.MustCompile(`\\*( feat\..*)`)
	songTitle = rmFeatDot.ReplaceAllLiteralString(songTitle, "")
	rmProdDot := regexp.MustCompile(`\\*( prod\..*)`)
	songTitle = rmProdDot.ReplaceAllLiteralString(songTitle, "")

	return strings.TrimSpace(songTitle)
}
func RmParens(songTitle string) string {
	rmParens := regexp.MustCompile(`\\*\([^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[[^)]*\]*`)
	songTitle = rmParens.ReplaceAllLiteralString(songTitle, "")
	songTitle = rmSquareBrackets.ReplaceAllLiteralString(songTitle, "")
	return strings.TrimSpace(songTitle)
}

func TrySingleArtist(artistName string) string {
	s := strings.Split(artistName, ", ")
	return s[0]
}
