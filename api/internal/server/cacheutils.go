package server

import (
	oxs "api/internal/oxsapi"
	"common/pb"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const (
	redisCacheTTL = time.Minute * 10
	redisPlayerPrefix   = "player:"
	redisLeaderboardKey = "leaderboard"
)

func getPlayerStat(ctx context.Context, ids []int64, rdb redis.UniversalClient) ([]oxs.PlayerStatJSON, error) {
	buf := make([]*oxs.PlayerStatJSON, len(ids))
	redisRes := make([]*pb.PlayerStats, len(ids))
	cached := 0
	for i, id := range ids {
		data, err := rdb.Get(ctx, redisPlayerPrefix+strconv.FormatInt(id, 10)).Bytes()
		if errors.Is(err, redis.Nil) {
			continue
		}
		cached++
		redisRes[i] = &pb.PlayerStats{}
		proto.Unmarshal(data, redisRes[i]) // with a lots of ids, it makes sense to make umnarshal queue and process it later
	}
	needToReq := []int64{}
	for i, id := range ids {
		if redisRes[i] == nil {
			needToReq = append(needToReq, id)
			continue
		}
		buf[i] = &oxs.PlayerStatJSON{
			PlayerId:     strconv.FormatInt(redisRes[i].GetId(), 10),
			MatchCount:   int(redisRes[i].GetMatchCount()),
			FavoriteHero: redisRes[i].GetFavoriteHero(),
			Places:       redisRes[i].GetPlaces(),
			Heroes:       make(map[string]oxs.PlayerStatHeroJSON),
		}
		for _, m := range redisRes[i].Matches {
			buf[i].Matches = append(buf[i].Matches, oxs.PlayerStatMatchJSON{
				HeroName:     m.GetHeroName(),
				MainTalent:   m.GetMainTalent(),
				Items:        m.GetItems(),
				Kills:        int(m.GetKills()),
				Deaths:       int(m.GetDeaths()),
				Place:        int(m.GetPlace()),
				RatingChange: int(m.GetRatingChange()),
				EndTime:      int(m.GetEndTime()),
				Date:         uintptr(m.GetDate()),
			})
		}
		for k, v := range redisRes[i].Heroes {
			buf[i].Heroes[k] = oxs.PlayerStatHeroJSON{
				MatchCount: int(v.MatchCount),
				Kills:      v.GetKills(),
				Deaths:     v.GetDeaths(),
				Rating:     v.GetRating(),
				Places:     v.GetPlaces(),
			}
		}
	}
	stats, err := oxs.GetPlayerStat(ctx, needToReq) // +-1s
	if err != nil {
		return nil, err
	}
	if cached+len(stats) != len(ids) {
		return nil, errors.New("something wrong with api or db")
	}
	res := make([]oxs.PlayerStatJSON, len(ids))
	for i, statsIdx := 0, 0; i < len(res); i++ {
		if buf[i] != nil {
			res[i] = *buf[i]
			continue
		}
		res[i] = stats[statsIdx]
		statsIdx++
	}
	go func() {
		for _, s := range stats {
			id, _ := strconv.ParseInt(s.PlayerId, 10, 64)
			ps := &pb.PlayerStats{
				Id:           id,
				MatchCount:   uint32(s.MatchCount),
				FavoriteHero: s.FavoriteHero,
				Places:       s.Places,
				Heroes:       make(map[string]*pb.PlayerStats_HeroStat),
			}
			for _, m := range s.Matches {
				ps.Matches = append(ps.Matches, &pb.PlayerStats_Match{
					HeroName:     m.HeroName,
					MainTalent:   m.MainTalent,
					Items:        m.Items,
					Kills:        uint32(m.Kills),
					Deaths:       uint32(m.Deaths),
					Place:        uint32(m.Place),
					RatingChange: int32(m.RatingChange),
					EndTime:      uint64(m.EndTime),
					Date:         uint64(m.Date),
				})
			}
			for k, v := range s.Heroes {
				ps.Heroes[k] = &pb.PlayerStats_HeroStat{
					MatchCount: int32(v.MatchCount),
					Kills:      v.Kills,
					Deaths:     v.Deaths,
					Rating:     v.Rating,
					Places:     v.Places,
				}
			}
			encoded, err := proto.Marshal(ps)
			if err != nil {
				continue
			}
			rdb.Set(context.Background(), redisPlayerPrefix+s.PlayerId, encoded, redisCacheTTL)
		}
	}()
	return res, nil
}

func getMatchLeaderboard(ctx context.Context, rdb redis.UniversalClient) ([]oxs.LeaderboardJSON, error) {
	data, err := rdb.Get(ctx, redisLeaderboardKey).Bytes()
	if errors.Is(err, redis.Nil) {
		res, err := oxs.GetMatchLeaderboard(ctx)
		cacheLb := &pb.Leaderboard{}
		for _, p := range res {
			cacheLb.Players = append(cacheLb.Players, &pb.Leaderboard_Player{
				Id: p.PlayerId,
				MatchCount: uint32(p.MatchCount),
				FavoriteHero: p.FavoriteHero,
				Rating: uint32(p.Rating),
			})
		}
		encoded, protoErr := proto.Marshal(cacheLb)
		if protoErr == nil {
			go func(){
				rdb.Set(context.Background(), redisLeaderboardKey, encoded, redisCacheTTL)
			}()
		}
		return res, err
	}
	lb := &pb.Leaderboard{}
	err = proto.Unmarshal(data, lb)
	if err != nil {
		return nil, err
	}
	res := make([]oxs.LeaderboardJSON, 100)
	for i, p := range lb.GetPlayers() {
		res[i] = oxs.LeaderboardJSON{
			PlayerId:     p.GetId(),
			FavoriteHero: p.GetFavoriteHero(),
			MatchCount:   int(p.GetMatchCount()),
			Rating:       int(p.GetRating()),
		}
	}
	return res, nil
}
