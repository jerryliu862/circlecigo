package google

import (
	"17live_wso_be/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) AuthGoogleUser(ctx context.Context, token string) (model.GoogleUser, error) {
	user, err := c.fetchGoogleUser(ctx, token)
	if err != nil {
		log.Errorf("auth google user fail: %s", err.Error())
		return user, err
	}

	return user, nil
}

func (c *Client) fetchGoogleUser(ctx context.Context, token string) (model.GoogleUser, error) {
	var user model.GoogleUser

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return user, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	q := req.URL.Query()
	q.Add(c.TokenQuery, token)
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return user, err
	} else if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &user); err != nil {
		return user, err
	}

	return user, nil
}
