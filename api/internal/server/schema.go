package server

type ResponseSchema struct {
	Data   any            `json:"data,omitempty"`
	Errors *[]ErrorSchema `json:"errors,omitempty"`
}

type ErrorCode int

const (
	UnspecifiedError ErrorCode = iota
	InternalError
	ParseError
	IncorrectParametersError
	ExternalApiError
)

type ErrorSchema struct {
	Code ErrorCode `json:"code,omitempty"`
	Desc string    `json:"desc,omitempty"`
}

type ReviewRequest struct {
	PairIdSchema
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

type PairIdSchema struct {
	First  int64 `json:"first_id,omitempty"`
	Second int64 `json:"second_id,omitempty"`
}

type PairFloatSchema struct {
	First  float32 `json:"first"`
	Second float32 `json:"second"`
}

type FriendsResponse struct {
	PairIdSchema
	GamesInSameLobby   int             `json:"games_in_same_lobby"`
	AvgPlaces          PairFloatSchema `json:"avg_places"`
	SameLobbyAvgPlaces PairFloatSchema `json:"same_lobby_avg_places"`
	PtsGained          PairFloatSchema `json:"pts_gained"`
}

type SmurfResponse struct {
	PairIdSchema
	// unix milli time
	LastMatchTogether int64 `json:"last_match_together_time"`
	HeroPools         struct {
		First  map[string]float32 `json:"first"`
		Second map[string]float32 `json:"second"`
	} `json:"hero_pools"`
}
