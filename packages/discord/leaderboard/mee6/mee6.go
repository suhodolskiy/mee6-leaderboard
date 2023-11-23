package mee6

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetLeaderboard(ctx context.Context, guildID string) (*Leaderboard, error) {
	url := fmt.Sprintf("https://mee6.xyz/api/plugins/levels/leaderboard/%s", guildID)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var leaderboard Leaderboard

		if err := json.NewDecoder(resp.Body).Decode(&leaderboard); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}

		return &leaderboard, nil
	case http.StatusNotFound:
		return nil, ErrGuildNotFound
	default:
		return nil, fmt.Errorf("get leader board: %d", resp.StatusCode)
	}
}
