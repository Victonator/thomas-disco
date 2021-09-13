package music

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		AddedBy struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			ID   string `json:"id"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"added_by"`
		IsLocal      bool        `json:"is_local"`
		PrimaryColor interface{} `json:"primary_color"`
		SharingInfo  struct {
			ShareID  string `json:"share_id"`
			ShareURL string `json:"share_url"`
			URI      string `json:"uri"`
		} `json:"sharing_info"`
		Track struct {
			Album struct {
				AlbumType string `json:"album_type"`
				Artists   []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					Height int    `json:"height"`
					URL    string `json:"url"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				TotalTracks          int    `json:"total_tracks"`
				Type                 string `json:"type"`
				URI                  string `json:"uri"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			DiscNumber  int  `json:"disc_number"`
			DurationMs  int  `json:"duration_ms"`
			Episode     bool `json:"episode"`
			Explicit    bool `json:"explicit"`
			ExternalIds struct {
				Isrc string `json:"isrc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href        string `json:"href"`
			ID          string `json:"id"`
			IsLocal     bool   `json:"is_local"`
			IsPlayable  bool   `json:"is_playable"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			Track       bool   `json:"track"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
		} `json:"track"`
		VideoThumbnail struct {
			URL interface{} `json:"url"`
		} `json:"video_thumbnail"`
	} `json:"items"`
}

//https://api.spotify.com/v1/playlists/2E0Yp9EH27aq68HRhvDXWk/tracks?offset=0&limit=100&additional_types=track%2Cepisode&market=BE

func (m *MusicCommand) SpotifyToSearch(id string) []string {
	tracks := []string{}

	hasMore := true
	offset := 0
	for hasMore {
		// call spotify
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?offset=%d&limit=100&additional_types=track&market=BE", id, offset), nil)
		if err != nil {
			log.Printf("error: %s", err)
			return tracks
		}
		req.Header.Set("Authorization", "Bearer "+m.SpotifyToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error: %s", err)
			return tracks
		}
		defer resp.Body.Close()

		if err != nil {
			log.Println("Error getting spotify playlist", "err", err)
			return tracks

		}

		body, _ := ioutil.ReadAll(resp.Body)

		data := Response{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Println("Error parsing spotify playlist", "err", err, string(body))
			return tracks
		}

		for _, item := range data.Items {
			if len(item.Track.Artists) == 0 {
				tracks = append(tracks, item.Track.Name)
			} else {
				tracks = append(tracks, fmt.Sprintf("%s %s", item.Track.Name, item.Track.Artists[0].Name))
			}
		}

		if len(data.Items) < 100 {
			hasMore = false
		} else {
			offset += 100
		}
	}

	return tracks
}
