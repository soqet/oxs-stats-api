package analyzer

import (
	api "api/internal/oxsapi"
)

type FriendReview struct {
	first            api.PlayerStatJSON
	second           api.PlayerStatJSON
	GamesInSameLobby int
	AllAvgPlaces     AvgPlaces
	PartyAvgPlaces   AvgPlaces
	PtsGained        struct {
		First  int
		Second int
	}
}

type AvgPlaces struct {
	First  float32
	Second float32
}

func DetectFriends(first, second api.PlayerStatJSON) FriendReview {
	res := FriendReview{}
	res.first = first
	res.second = second
	res.fillAvgPlaces()
	return res
}

func (fr *FriendReview) fillAvgPlaces() {
	matchesChecked := 0
	firstRankedMatches := 0
	secondRankedMatches := 0
	for j := matchesChecked; j < len(fr.second.Matches); j++ {
		if fr.second.Matches[j].RatingChange == 0 {
			continue
		}
		fr.AllAvgPlaces.Second += float32(fr.second.Matches[j].Place)
		secondRankedMatches++
	}
	for i := 0; i < len(fr.first.Matches); i++ {
		if fr.first.Matches[i].RatingChange == 0 {
			continue
		}
		firstRankedMatches++
		fr.AllAvgPlaces.First += float32(fr.first.Matches[i].Place)
		for j := matchesChecked; j < len(fr.second.Matches); j++ {
			if fr.second.Matches[j].RatingChange == 0 {
				continue
			}
			if fr.first.Matches[i].GetTime().Equal(fr.second.Matches[j].GetTime()) {
				fr.GamesInSameLobby++
				fr.PartyAvgPlaces.First += float32(fr.first.Matches[i].Place)
				fr.PartyAvgPlaces.Second += float32(fr.second.Matches[j].Place)
				fr.PtsGained.First += fr.first.Matches[i].RatingChange
				fr.PtsGained.Second += fr.second.Matches[j].RatingChange
				matchesChecked = j + 1
				break
			}
		}
	}
	if firstRankedMatches != 0 {
		fr.AllAvgPlaces.First /= float32(firstRankedMatches)
	} else {
		fr.AllAvgPlaces.First = 0
	}
	if secondRankedMatches != 0 {
		fr.AllAvgPlaces.Second /= float32(secondRankedMatches)
	} else {
		fr.AllAvgPlaces.Second = 0
	}
	if fr.GamesInSameLobby != 0 {
		fr.PartyAvgPlaces.First /= float32(fr.GamesInSameLobby)
		fr.PartyAvgPlaces.Second /= float32(fr.GamesInSameLobby)
	} else {
		fr.PartyAvgPlaces.First = 0
		fr.PartyAvgPlaces.Second = 0
	}
}
