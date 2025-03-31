package util

// Structs for JSON parsing
type JWK struct {
	E   string `json:"e"`
	N   string `json:"n"`
	Kty string `json:"kty"`
	Kid string `json:"kid"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type ReqTokenValidation struct {
	Token string `json:"token" valid:"required"`
}

type GoogleJWTClaims struct {
	Audience        string `json:"aud"`
	AuthorizedParty string `json:"azp"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	Expiration      int64  `json:"exp"`
	FamilyName      string `json:"family_name"`
	GivenName       string `json:"given_name"`
	Issuer          string `json:"iss"`
	JWTID           string `json:"jti"`
	Name            string `json:"name"`
	NotBefore       int64  `json:"nbf"`
	Picture         string `json:"picture"`
	Subject         string `json:"sub"`
}
