package analyzer

import (
	api "api/internal/oxsapi"
	"time"
)

type SmurfReview struct {
	first                 api.PlayerStatJSON
	second                api.PlayerStatJSON
	LastMatchIntersection time.Time
	MatchesIntersection   []MatchIntersection
	HP                    HeroPools
	SameHeroNumber        int
}

type MatchIntersection struct {
	First  api.PlayerStatMatchJSON
	Second api.PlayerStatMatchJSON
}

type HeroPools struct {
	First  map[string]float32
	Second map[string]float32
}

func haveIntersection(start1, end1, start2, end2 time.Time) bool {
	return ((start1.After(start2) || start1.Equal(start2)) && start1.Before(end2)) ||
		(end1.After(start2) && (end1.Before(end2) || end1.Equal(end2)))
}

func DetectSmurf(first, second api.PlayerStatJSON) SmurfReview {
	res := SmurfReview{}
	res.first = first
	res.second = second
	res.fillMatchIntersection()
	res.fillHeroPool()
	return res
}

func (dr *SmurfReview) fillMatchIntersection() {
	var mi []MatchIntersection
	for i, j := 0, 0; i < len(dr.first.Matches) && j < len(dr.second.Matches); {
		curFirst := dr.first.Matches[i]
		curSecond := dr.second.Matches[j]
		start1 := curFirst.GetTime()
		end1 := start1.Add(time.Second * time.Duration(curFirst.EndTime))
		start2 := curSecond.GetTime()
		end2 := start2.Add(time.Second * time.Duration(curSecond.EndTime))
		if haveIntersection(start1, end1, start2, end2) {
			mi = append(mi, MatchIntersection{
				First:  curFirst,
				Second: curSecond,
			})
		}
		// because of inversion (first match in slice is the latest and)
		if end1.After(end2) {
			i++
		} else {
			j++
		}
	}
	dr.MatchesIntersection = mi
}

func (dr *SmurfReview) fillHeroPool() {
	firstPool := map[string]float32{}
	secondPool := map[string]float32{}
	for _, m := range dr.first.Matches {
		if p, ok := firstPool[m.HeroName]; ok {
			firstPool[m.HeroName] = p + 2
		}
	}
	for _, m := range dr.second.Matches {
		if p, ok := secondPool[m.HeroName]; ok {
			secondPool[m.HeroName] = p + 2
		}
	}
	k := 0
	for h := range firstPool {
		if _, ok := secondPool[h]; ok {
			k++
		}
	}
	dr.SameHeroNumber = k
	dr.HP = HeroPools{
		First:  firstPool,
		Second: secondPool,
	}
}
