package oxsapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LeaderboardJSON struct {
	PlayerId     int64  `json:"playerId"`
	Rating       int    `json:"rating"`
	FavoriteHero string `json:"favoriteHero"`
	MatchCount   int    `json:"matchCount"`
}

type PlayerStatJSON struct {
	PlayerId     string                        `json:"playerId"`
	MatchCount   int                           `json:"matchCount"`
	FavoriteHero string                        `json:"favoriteHero"`
	Places       map[string]int32              `json:"places"`
	Heroes       map[string]PlayerStatHeroJSON `json:"heroes"`
	Matches      []PlayerStatMatchJSON         `json:"matches"`
}

type PlayerStatHeroJSON struct {
	MatchCount int              `json:"matchCount"`
	Kills      string           `json:"kills"`
	Deaths     string           `json:"deaths"`
	Rating     string           `json:"rating"`
	Places     map[string]int32 `json:"places"` // 1 - 6
}

type PlayerStatMatchJSON struct {
	HeroName     string            `json:"heroName"`
	MainTalent   string            `json:"mainTalent"`
	Items        map[string]string `json:"items"` // 0 - 5
	Kills        int               `json:"kills"`
	Deaths       int               `json:"deaths"`
	Place        int               `json:"place"`
	RatingChange int               `json:"ratingChange"`
	EndTime      int               `json:"endTime"`
	Date         uintptr           `json:"date"`
}

type PublicLeaderboardJSON struct {
	AvatarUrl    string `json:"avatarUrl"`
	FavoriteHero string `json:"favoriteHero"`
	MatchCount   int    `json:"matchCount"`
	Nick         string `json:"nickname"`
	Rating       int    `json:"rating"`
	SteamUrl     string `json:"steamUrl"`
}

func GetMatchLeaderboard(ctx context.Context) ([]LeaderboardJSON, error) {
	url := "https://stats.dota1x6.com/api/match_leaderboard"
	method := http.MethodPost

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	input, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []LeaderboardJSON

	json.Unmarshal(input, &data)
	return data, nil
}

func GetPublicLeaderboard() ([]PublicLeaderboardJSON, error) {
	url := "https://stats.dota1x6.com/api/leaderboard"
	method := http.MethodGet

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	input, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []PublicLeaderboardJSON

	json.Unmarshal(input, &data)
	return data, nil
}

/*
strings.NewReader(func()string{
		res := strings.Builder{}
		res.WriteRune('[')
		for i, id := range steam32ids {
			res.WriteString(strconv.FormatInt(id, 10))
			if i < len(steam32ids) - 1 {
				res.WriteString(",")
			}
		}
		res.WriteRune(']')
		return res.String()
	}())
*/

func GetPlayerStat(ctx context.Context, steam32ids []int64) ([]PlayerStatJSON, error) {
	url := "https://stats.dota1x6.com/api/match_details"
	method := http.MethodPost
	client := &http.Client{}
	if len(steam32ids) <= 0 {
		return []PlayerStatJSON{}, nil
	}
	reqData, err := json.Marshal(&steam32ids)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, errors.New("bad request")
	}
	input, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []PlayerStatJSON
	// input = fixDetailsJSON(input)
	json.Unmarshal(input, &data)
	return data, nil
}

func (m *PlayerStatMatchJSON) GetTime() time.Time {
	return time.UnixMilli(int64(m.Date)).Add(time.Hour * -3)
}
