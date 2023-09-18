package midleware

import (
	"chat/pkg/server"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenString string) (jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			cfg, err := server.LoadConfig()
			if err != nil {
				return nil, fmt.Errorf("error loading configuration: %w", err)
			}
			return []byte(cfg.JwtKey), nil
		},
	)

	if err != nil || !token.Valid {
		fmt.Println("Not a valid token")
		return jwt.StandardClaims{}, errors.New("not a valid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		fmt.Println("Can't parse the claims")
		return jwt.StandardClaims{}, errors.New("can't parse the claims")
	}

	return *claims, nil
}

func AuthenticateUser(ctx *gin.Context) {
	authHelper(ctx, "user")
}

func authHelper(ctx *gin.Context, user string) {

	tokenString, err := ctx.Cookie(user + "-auth")
	if err != nil || tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Unauthorized, Please Login",
		})
		return
	}

	claims, err := ValidateToken(tokenString)
	if err != nil || tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Unauthorized, Please Login token not valid",
		})
		return
	}

	if time.Now().Unix() > claims.ExpiresAt {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Need Re-Login time expired",
		})
		return
	}
	ctx.Set("userId", fmt.Sprint(claims.Id))
	ctx.Set("email", fmt.Sprint(claims.Audience))

}
