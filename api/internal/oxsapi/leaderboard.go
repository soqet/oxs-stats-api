package oxsapi

import (
	"bytes"
	"context"
	"regexp"
	"strconv"
)

type Stat struct {
	HeroName   string
	TalentName string
	AvgPlace   float32
	MatchCount int
}

func GetStatsLeaderboard(ctx context.Context) ([]StatJSON, error) {
	talents := map[string]Stat{} // key : hero + "-" + talent
	lb, err := GetMatchLeaderboard(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range lb {
		psArr, err := GetPlayerStat(context.Background(), []int64{l.PlayerId})
		if err != nil || len(psArr) < 1 {
			return nil, err
		}
		ps := psArr[0]
		for _, m := range ps.Matches {
			if m.MainTalent == "" {
				// fmt.Println(m, l.PlayerId)
				continue
			}
			t, ok := talents[m.HeroName+"-"+m.MainTalent]
			if !ok {
				t = Stat{
					HeroName:   m.HeroName,
					TalentName: m.MainTalent,
					AvgPlace:   float32(m.Place),
					MatchCount: 1,
				}
			} else {
				t.AvgPlace = (t.AvgPlace*float32(t.MatchCount) + float32(m.Place)) / float32(t.MatchCount+1)
				t.MatchCount++
			}
			talents[m.HeroName+"-"+m.MainTalent] = t
		}
	}
	s := []StatJSON{}
	for _, v := range talents {
		s = append(s, StatJSON{
			HeroName:   v.HeroName,
			TalentName: v.TalentName,
			AvgPlace:   strconv.FormatFloat(float64(v.AvgPlace), 'f', 2, 64),
			MatchCount: uintptr(v.MatchCount),
		})
	}
	return s, nil
}

// пиздец костыль, но это не моя вина что у ксено блять кривые жсоны приходят
func fixDetailsJSON(data []byte) []byte {
	var re = regexp.MustCompile(`(?mU)"items"\s*:\s*\[.*\]`)
	return re.ReplaceAllFunc(data, replDetailsJSON)
}

func replDetailsJSON(arr []byte) []byte {
	arr = bytes.ReplaceAll(arr, []byte{'['}, []byte{'{'})
	arr = bytes.ReplaceAll(arr, []byte{']'}, []byte{'}'})
	var re = regexp.MustCompile(`(?mU)".*"`)
	bracket := bytes.Index(arr, []byte("{"))
	arr = append(arr[:bracket],
		re.ReplaceAllFunc(arr[bracket:], func(b []byte) []byte {
			return append([]byte(`"|":`), b...)
		})...,
	)
	itemIdx := byte('0')
	for i := range arr {
		if arr[i] == '|' {
			arr[i] = itemIdx
			itemIdx++
		}
	}
	return arr
}
