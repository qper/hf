package domain

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User          User     `json:"user"`
	RecoveryCodes []string `json:"recovery_codes"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RecoverRequest struct {
	Username     string `json:"username"`
	RecoveryCode string `json:"recovery_code"`
}

type RecoverResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type RecoveryCodeStatusResponse struct {
	Remaining int `json:"remaining"`
}

type RecoveryCodeRegenerationRequest struct {
	Password string `json:"password"`
}

type RecoveryCodeRegenerationResponse struct {
	RecoveryCodes []string `json:"recovery_codes"`
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
