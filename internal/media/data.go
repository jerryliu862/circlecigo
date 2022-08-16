package media

import (
	"17live_wso_be/internal/model"
	"17live_wso_be/util"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) FetchAccessToken(ctx context.Context) (string, error) {
	var token string

	var auth model.MediaApiAuthentication

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	body := strings.NewReader(params.Encode())

	url := c.Authentication.Server + c.Authentication.Path

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return token, err
	}
	req.SetBasicAuth(c.Authentication.ClientID, c.Authentication.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	} else if resp.StatusCode != http.StatusOK {
		return token, fmt.Errorf("unexpected status code when get access token: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &auth); err != nil {
		return token, err
	}

	token = fmt.Sprintf("%s %s", auth.TokenType, auth.AccessToken)

	return token, err
}

func (c *Client) FetchCampaignData(ctx context.Context, skipList []int) ([]model.MediaApiCampaignSet, error) {
	log.Infof("fetch campaign data with skip sync list: %v", skipList)

	var data []model.MediaApiCampaignSet

	httpClient := &http.Client{}

	url := c.Campaign.Server + c.Campaign.Path

	offset := 0
	limit, err := strconv.Atoi(c.Campaign.LimitValue)
	if err != nil {
		log.Errorf("invalid campaign limit value in config: %s", c.Campaign.LimitValue)
		return data, err
	}

	for offset >= 0 {
		var d []model.MediaApiCampaignSet

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return data, err
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set(c.Campaign.TokenHeaderKey, c.Campaign.Token)

		q := req.URL.Query()
		q.Add(c.Campaign.StatusQuery, c.Campaign.StatusValue)
		q.Add("offset", strconv.Itoa(offset))
		q.Add(c.Campaign.LimitQuery, c.Campaign.LimitValue)
		req.URL.RawQuery = q.Encode()

		resp, err := httpClient.Do(req)
		if err != nil {
			return data, err
		} else if resp.StatusCode != http.StatusOK {
			return data, fmt.Errorf("unexpected status code when sync campaign: %d", resp.StatusCode)
		}

		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(respBody, &d); err != nil {
			return data, err
		}

		offset += limit

		if len(d) < limit {
			log.Infof("reach last page with limit %d and offset %d", limit, offset-limit)
			offset = -1
		}

		for _, v := range d {
			if util.ContainInt(skipList, v.Campaign.Id) {
				continue
			}
			data = append(data, v)
		}
	}

	return data, nil
}

func (c *Client) FetchLeaderboardData(ctx context.Context, leaderboardID string) (model.MediaApiLeaderboardDetail, error) {
	var data model.MediaApiLeaderboardDetail

	httpClient := &http.Client{}

	url := c.Leaderboard.Server + c.Leaderboard.Path
	url = strings.Replace(url, "?", leaderboardID, 1)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set(c.Leaderboard.TokenHeaderKey, c.Leaderboard.Token)

	q := req.URL.Query()
	q.Add(c.Leaderboard.CursorQuery, c.Leaderboard.CursorValue)
	q.Add(c.Leaderboard.CountQuery, c.Leaderboard.CountValue)
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return data, err
	} else if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("unexpected status code when sync campaign leaderboard: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *Client) FetchStreamerData(ctx context.Context, streamerID string, token string) (model.MediaApiStreamer, error) {
	var data model.MediaApiStreamer

	httpClient := &http.Client{}

	url := c.Streamer.Server + c.Streamer.Path
	url = strings.Replace(url, "?", streamerID, 1)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set(c.Streamer.TokenHeaderKey, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return data, err
	} else if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("unexpected status code when sync streamer: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *Client) FetchStreamerContractData(ctx context.Context, streamerID, token string) (model.MediaApiStreamerContract, error) {
	var data model.MediaApiStreamerContract

	httpClient := &http.Client{}

	url := c.StreamerContract.Server + c.StreamerContract.Path
	url = strings.Replace(url, "?", streamerID, 1)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set(c.Streamer.TokenHeaderKey, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return data, err
	} else if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("unexpected status code when sync streamer contract: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &data); err != nil {
		return data, err
	}

	return data, nil
}
