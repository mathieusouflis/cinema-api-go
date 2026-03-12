package jwt

type JWTClaims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}
