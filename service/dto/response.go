package dto

type DefaultResponse struct {
	Message string `json:"message,omitempty"`
}

type LoginResponse struct {
	IDToken      string `json:"idToken"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}
