package spotifyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// Amount of time to wait to retry if we are told to retry
	// but not provided a wait-interval
	defualtRetryDuration = time.Second * 5

	// rateLimitExceededStatusCode is the code that the server returns when our
	// request frequency is too high.
	rateLimitExceededStatusCode = 429

	spotifyBaseURL = "https://api.spotify.com/v1/"

	// SpotifyTimeRangeShort is ~4 weeks
	SpotifyTimeRangeShort = "short_term"

	// SpotifyTimeRangeMedium is ~6 months
	SpotifyTimeRangeMedium = "medium_term"

	// SpotifyTimeRangeLong is several years
	SpotifyTimeRangeLong = "long_term"
)

// Client stores the client information for working with the spotify web API.
type Client struct {
	client    *http.Client
	baseURL   string
	AutoRetry bool
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

// retryDuration contains the logic for determining when to retry request.
func retryDuration(resp *http.Response) time.Duration {
	str := resp.Header.Get("Retry-After")
	if str == "" {
		return defualtRetryDuration
	}
	numSeconds, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return defualtRetryDuration
	}
	return time.Duration(numSeconds) * time.Second
}

// get performs a http get request on the url provided
func (c *Client) get(requestURL string, result interface{}) error {
	for {
		resp, err := c.client.Get(requestURL)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}

		if resp.StatusCode == http.StatusNoContent {
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("spotify returned a status code of '%d' with a status text of '%s'", resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		err = json.NewDecoder(resp.Body).Decode(result)

		if err != nil {
			return err
		}

		break
	}

	return nil
}

// GetUserTopTracks returns a PagingTrack object containing the users top n listed to
// tracks in the time_range where n is the limit.  For more info on timeRange see
// https://developer.spotify.com/documentation/web-api/reference/personalization/get-users-top-artists-and-tracks/
func (c *Client) GetUserTopTracks(limit int, timeRange string) (*PagingTrack, error) {
	requestURL, err := url.Parse(fmt.Sprintf("%sme/top/tracks", spotifyBaseURL))
	if err != nil {
		return nil, fmt.Errorf("Malformed URL: %s", err.Error())
	}
	// Prepare Query Parameters
	params := url.Values{}
	params.Add("time_range", timeRange)
	params.Add("limit", strconv.Itoa(limit))

	// Add Query Parameters to the URL
	requestURL.RawQuery = params.Encode() // Escape Query Parameters

	var tracks PagingTrack
	err = c.get(requestURL.String(), &tracks)
	if err != nil {
		return nil, err
	}

	return &tracks, nil
}
