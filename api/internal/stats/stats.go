package stats

import (
	oxs "api/internal/oxsapi"
	"strings"
)

// type HeroTalentStringed string

type HeroTalent struct {
	HeroName   string `json:"hero_name"`
	TalentName string `json:"talent_name"`
}

const delim = "+"

func (ht HeroTalent) String() string {
	b := strings.Builder{}
	b.WriteString(ht.HeroName)
	b.WriteString(delim)
	b.WriteString(ht.TalentName)
	return b.String()
}

func NewHeroTalent(str string) HeroTalent {
	s := strings.Split(string(str), delim)
	if len(s) < 1 {
		return HeroTalent{}
	} else if len(s) < 2 {
		return HeroTalent{
			HeroName: s[0],
		}
	}
	return HeroTalent{
		HeroName:   s[0],
		TalentName: s[1],
	}
}

type HeroStats struct {
	HeroTalent
	Pickrate   float32     `json:"pickrate"`
	AvgPlace   float32     `json:"avg_place"`
	Matchcount int         `json:"match_count"`
	Places     map[int]int `json:"places"`
	AvgPts     float32     `json:"avg_pts"` // avg pts earned as this hero
	// Counters      map[string]float32 // avg places vs hero+talent
	// gamesCounters map[string]int
}

type StatReport struct {
	Stats map[string]HeroStats
}

func (sr StatReport) fillCounters(matches []oxs.PlayerStatMatchJSON) {
	games := map[uintptr][]oxs.PlayerStatMatchJSON{}
	for _, m := range matches {
		if m.RatingChange == 0 {
			continue
		}
		games[m.Date] = append(games[m.Date], m)
	}
	for _, v := range games {
		if len(v) <= 1 {
			continue
		}
		placings := map[string]int{}
		for _, g := range v {
			placings[HeroTalent{
				HeroName:   g.HeroName,
				TalentName: g.MainTalent,
			}.String()] = g.Place
		}
		// for ht, g := range placings {
		// 	sr.Stats[ht].gamesCounters[]++

		// }
	}
}

func MakeStats(matches []oxs.PlayerStatMatchJSON) StatReport {
	sr := StatReport{
		Stats: map[string]HeroStats{},
	}
	rankedGames := 0
	for _, m := range matches {
		if m.RatingChange == 0 {
			continue
		}
		rankedGames++
		hashName := HeroTalent{
			HeroName:   m.HeroName,
			TalentName: m.MainTalent,
		}.String()
		var stats HeroStats
		var ok bool
		if stats, ok = sr.Stats[hashName]; !ok {
			stats.HeroName = m.HeroName
			stats.TalentName = m.MainTalent
		}
		if stats.Places == nil {
			stats.Places = map[int]int{}
		}
		stats.AvgPlace = (stats.AvgPlace*float32(stats.Matchcount) + float32(m.Place)) / float32(stats.Matchcount+1)
		stats.Matchcount++
		stats.Places[m.Place]++
		sr.Stats[hashName] = stats
	}
	for k := range sr.Stats {
		v := sr.Stats[k]
		v.AvgPts = 40*(float32(v.Places[1])/float32(v.Matchcount)) +
			30*(float32(v.Places[2])/float32(v.Matchcount)) +
			10*(float32(v.Places[3])/float32(v.Matchcount)) -
			10*(float32(v.Places[4])/float32(v.Matchcount)) -
			30*(float32(v.Places[5])/float32(v.Matchcount)) -
			40*(float32(v.Places[6])/float32(v.Matchcount))
		v.Pickrate = float32(v.Matchcount) / float32(rankedGames)
		sr.Stats[k] = v
	}
	return sr
}
