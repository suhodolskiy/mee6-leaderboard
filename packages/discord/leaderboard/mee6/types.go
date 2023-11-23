package mee6

type Leaderboard struct {
	Players []struct {
		Avatar               string `json:"avatar"`
		DetailedXp           []int  `json:"detailed_xp"`
		Discriminator        string `json:"discriminator"`
		GuildID              string `json:"guild_id"`
		ID                   string `json:"id"`
		IsMonetizeSubscriber bool   `json:"is_monetize_subscriber"`
		Level                int    `json:"level"`
		MessageCount         int    `json:"message_count"`
		MonetizeXpBoost      int    `json:"monetize_xp_boost"`
		Username             string `json:"username"`
		Xp                   int    `json:"xp"`
	} `json:"players"`
}

type Error struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	StatusCode int `json:"status_code"`
}
