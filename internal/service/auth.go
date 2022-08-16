package service

import (
	"17live_wso_be/config"
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserTokenClaims struct {
	Id int
	jwt.StandardClaims
}

func (c *Client) GetUserToken(ctx context.Context, uid int) (*model.UserToken, error) {
	now := time.Now()
	claims := new(UserTokenClaims)
	claims.Id = uid
	claims.Issuer = config.New().Jwt.Issure
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(config.New().Jwt.ExpiredHour * time.Hour).Unix()
	idToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := idToken.SignedString([]byte(config.New().Jwt.Hmac))

	if err != nil {
		log.Errorf("fail to issue user token: %s", err.Error())
		return nil, customError.New(customError.UnknownError)
	}

	return &model.UserToken{
		Token: token,
	}, nil
}

func (c *Client) ValidUserToken(ctx context.Context, userToken string) (*UserTokenClaims, error) {
	token, err := jwt.ParseWithClaims(userToken, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.New().Jwt.Hmac), nil
	})

	if err != nil {
		log.Errorf("invalid user token: %s", err.Error())
		return nil, customError.New(customError.InvalidUserToken)
	}

	if claims, ok := token.Claims.(*UserTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, customError.New(customError.InvalidUserToken)
}

func (c *Client) PermissionCheck(ctx context.Context, uid int, region string, authType int, minAuthLevel int) bool {
	level, err := c.RepositoryClient.GetUserAuthLevel(ctx, uid, region, authType)

	if err != nil {
		return false
	}

	return level >= minAuthLevel
}
