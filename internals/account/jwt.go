package account

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

var accessDuration int
var refreshDuration int
var jwtKey string

var refreshBaseDuration = time.Hour * 24
var accessBaseDuration = time.Second

func SetTokenizeConfig(
	accDur, rfrDur int,
	jwtK string,
) {
	accessDuration = accDur
	refreshDuration = rfrDur
	jwtKey = jwtK
}

type claimResponse struct {
	Id    ulid.ULID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

type OauthToken struct {
	UserId      ulid.ULID `json:"user_id"`
	DefaultBook ulid.ULID `json:"default_book"`
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	ExpiredAt   time.Time `json:"exp_time"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
	ExpiredAt    string `json:"expired_at"`
	Scope        string `json:"scope"`
}

func generateJwtToken(
	ctx context.Context,
	id ulid.ULID,
	defaultBook ulid.ULID,
	email string,
) (
	res LoginResponse,
	err error,
) {

	// 15 minute
	expirationAccessTime := time.Now().Add(time.Duration(accessDuration) * accessBaseDuration)
	accessClaims := claimResponse{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationAccessTime),
		},
	}
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	)
	accessTokenString, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return
	}
	accessDataToken := OauthToken{
		UserId:      id,
		Email:       email,
		DefaultBook: defaultBook,
		Token:       accessTokenString,
		ExpiredAt:   expirationAccessTime,
	}

	// one month
	refreshDurationTime := time.Now().Add(time.Duration(refreshDuration) * refreshBaseDuration)
	refreshClaims := claimResponse{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshDurationTime),
		},
	}

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	)

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		return
	}

	refreshDataToken := OauthToken{
		UserId:      id,
		Email:       email,
		DefaultBook: defaultBook,
		Token:       refreshTokenString,
		ExpiredAt:   refreshDurationTime,
	}

	res = LoginResponse{
		AccessToken:  accessDataToken.Token,
		RefreshToken: refreshDataToken.Token,
		Type:         "Bearer",
		ExpiredAt:    refreshDurationTime.Format(time.RFC3339),
		Scope:        "*",
	}

	return
}
