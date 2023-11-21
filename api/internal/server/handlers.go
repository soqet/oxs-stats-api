package server

import (
	"crypto/ecdsa"
	"net/http"
	"strconv"
	"time"

	"api/internal/analyzer"
	oxs "api/internal/oxsapi"
	"api/internal/stats"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func handleAuth(logger zerolog.Logger, privKey *ecdsa.PrivateKey) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

func handleRegister(logger zerolog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

func handleFriendReview(logger zerolog.Logger, rdb redis.UniversalClient) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := ReviewRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseSchema{Errors: &[]ErrorSchema{{Code: ParseError, Desc: "invalid payload"}}})
			return
		}
		playerStats, err := getPlayerStat(r.Context(), []int64{req.First, req.Second}, rdb)
		if err != nil || len(playerStats) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseSchema{Errors: &[]ErrorSchema{{Code: ExternalApiError, Desc: "1x6 api is unavaible"}}})
			return
		}
		firstStat := playerStats[0]
		secondStat := playerStats[1]
		rev := analyzer.DetectFriends(firstStat, secondStat)
		id1, _ := strconv.ParseInt(firstStat.PlayerId, 10, 64)
		id2, _ := strconv.ParseInt(secondStat.PlayerId, 10, 64)
		resp := ResponseSchema{
			Data: FriendsResponse{
				PairIdSchema: PairIdSchema{
					First:  id1,
					Second: id2,
				},
				GamesInSameLobby: rev.GamesInSameLobby,
				AvgPlaces: PairFloatSchema{
					First:  rev.AllAvgPlaces.First,
					Second: rev.AllAvgPlaces.Second,
				},
				SameLobbyAvgPlaces: PairFloatSchema{
					First:  rev.PartyAvgPlaces.First,
					Second: rev.PartyAvgPlaces.Second,
				},
				PtsGained: PairFloatSchema{
					First:  float32(rev.PtsGained.First),
					Second: float32(rev.PtsGained.Second),
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
}

func handleSmurfReview(logger zerolog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
}

func handleStats(logger zerolog.Logger, rdb redis.UniversalClient) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lb, err := getMatchLeaderboard(r.Context(), rdb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseSchema{
				Errors: &[]ErrorSchema{
					{
						Code: ExternalApiError,
						Desc: "1x6 api is unavaible",
					},
				},
			})
			return
		}
		ids := []int64{}
		for _, p := range lb {
			ids = append(ids, p.PlayerId)
		}
		players, err := getPlayerStat(r.Context(), ids, rdb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseSchema{
				Errors: &[]ErrorSchema{
					{
						Code: ExternalApiError,
						Desc: "1x6 api is unavaible",
					},
				},
			})
		}
		matches := []oxs.PlayerStatMatchJSON{}
		for _, p := range players {
			for _, m := range p.Matches {
				//* date of last patch
				if m.GetTime().Before(time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC)) {
					break
				}
				if m.RatingChange != 0 {
					matches = append(matches, m)
				}
			}
		}
		json.NewEncoder(w).Encode(ResponseSchema{
			Data: stats.MakeStats(matches).Stats,
		})
	})
}
