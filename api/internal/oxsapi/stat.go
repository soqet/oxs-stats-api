package oxsapi

import (
	"io"
	"net/http"
	"strconv"
)

type StatJSON struct {
	HeroName   string             `json:"heroName"`
	TalentName string             `json:"talentName"`
	Pickrate   string             `json:"pickrate"`
	AvgPlace   string             `json:"averagePlace"`
	MatchCount uintptr            `json:"matchCount"`
	Places     map[string]float64 `json:"places"`
	Probs      map[int]float64
	Ex         float64
	Dx         float64
	PtsAvg     float64
}

func (s *StatJSON) setProbs() {
	s.Probs = map[int]float64{}
	var sum float64
	for _, v := range s.Places {
		sum += v
	}
	for k, v := range s.Places {
		p, _ := strconv.Atoi(k)
		s.Probs[p] = v / sum
	}
}

func (s *StatJSON) SetStats() {
	s.setProbs()
	s.Ex = s.getEx()
	s.Dx = s.getDx()
	s.PtsAvg = s.getPtsAvg()
}

func (s StatJSON) getEx() (e float64) {
	for k, v := range s.Probs {
		e += float64(k) * v
	}
	return
}

func (s StatJSON) getDx() (d float64) {
	for k, v := range s.Probs {
		d += float64(k) * float64(k) * v
	}
	e := s.getEx()
	d -= e * e
	return
}

func (s StatJSON) getPtsAvg() (avg float64) {
	avg = 40*s.Probs[1] + 30*s.Probs[2] + 10*s.Probs[3] - 10*s.Probs[4] - 20*s.Probs[5] - 40*s.Probs[6]
	return avg
}

func getStatsPlus() ([]StatJSON, error) {
	url := "https://stats.dota1x6.com/api/heroes?talents=yes&period=four_days&pool=all"
	method := http.MethodGet

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("fallbackcookie", "token=token") // TOKEN HERE

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	input, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []StatJSON

	json.Unmarshal(input, &data)
	for i := 0; i < len(data); i++ {
		data[i].SetStats()
	}
	return data, nil
}
