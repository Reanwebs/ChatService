package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Middleware struct {
	Methods MiddlewareMethods
}

type MiddlewareMethods interface {
	ValidateToken(string) (jwt.StandardClaims, error)
	AuthenticateUser(*gin.Context)
	AuthHelper(*gin.Context, string)
}

func (m Middleware) ValidateToken(tokenString string) (jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("P6zqwuYDJomBGnleDYtF3pMyoN3sVaiy2BTbUNd566g"), nil
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

func (m Middleware) AuthenticateUser(ctx *gin.Context) {
	m.AuthHelper(ctx, "user")
}

func (m Middleware) AuthHelper(ctx *gin.Context, user string) {
	tokenString, err := ctx.Cookie(user + "-auth")
	if err != nil || tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"StatusCode": 401,
			"msg":        "Unauthorized, Please Login",
		})
		return
	}

	claims, err := m.ValidateToken(tokenString)
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
