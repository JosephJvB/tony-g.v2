package youtube

import "testing"

func TestYoutubeUtil(t *testing.T) {
	t.Run("GetYoutubeVideoID handles invalid url", func(t *testing.T) {
		id := GetYoutubeVideoID("invalid-url")
		if id != "" {
			t.Errorf("expected empty string, got %q", id)
		}
	})

	t.Run("GetYoutubeVideoID handles bandcamp url", func(t *testing.T) {
		id := GetYoutubeVideoID("https://jamilawoods.bandcamp.com/track/teach-me")
		if id != "" {
			t.Errorf("expected empty string, got %q", id)
		}
	})

	t.Run("GetYoutubeVideoID handles YouTube playlist URL", func(t *testing.T) {
		id := GetYoutubeVideoID("https://www.youtube.com/playlist?list=OLAK5uy_kd47NwWgbiLjeSiN99m424aoyarXAw-PE")
		if id != "" {
			t.Errorf("expected empty string, got %q", id)
		}
	})

	t.Run("GetYoutubeVideoID handles YouTube URL with v parameter", func(t *testing.T) {
		id := GetYoutubeVideoID("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
		if id != "dQw4w9WgXcQ" {
			t.Errorf("expected video ID 'dQw4w9WgXcQ', got %q", id)
		}
	})

	t.Run("GetYoutubeVideoID handles YouTube URL with v parameter and additional query params", func(t *testing.T) {
		id := GetYoutubeVideoID("https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley")
		if id != "dQw4w9WgXcQ" {
			t.Errorf("expected video ID 'dQw4w9WgXcQ', got %q", id)
		}
	})

	t.Run("GetYoutubeVideoId handles music.youtube domain", func(t *testing.T) {
		id := GetYoutubeVideoID("https://music.youtube.com/watch?v=dQw4w9WgXcQ")
		if id != "dQw4w9WgXcQ" {
			t.Errorf("expected video ID 'dQw4w9WgXcQ', got %q", id)
		}
	})

	t.Run("GetYoutubeVideoId handles youtu.be short URL", func(t *testing.T) {
		id := GetYoutubeVideoID("https://youtu.be/dQw4w9WgXcQ")
		if id != "dQw4w9WgXcQ" {
			t.Errorf("expected video ID 'dQw4w9WgXcQ', got %q", id)
		}
	})

	t.Run("GetYoutubeVideoId handles youtu.be short URL with query params", func(t *testing.T) {
		id := GetYoutubeVideoID("https://youtu.be/dQw4w9WgXcQ?param=value")
		if id != "dQw4w9WgXcQ" {
			t.Errorf("expected video ID 'dQw4w9WgXcQ', got %q", id)
		}
	})
}
