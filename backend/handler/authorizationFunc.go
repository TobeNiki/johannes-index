package handler

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/TobeNiki/Index/backend/database"
)

const (
	IdentityKey = "id"
	IndexKey    = "index"
)

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*database.M_User); ok {
		return jwt.MapClaims{
			IdentityKey: v.UserID,
			IndexKey:    v.ESIndexName,
		}
	}
	return jwt.MapClaims{}
}

func identityHandlerFunc(ctx *gin.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	return &database.M_User{
		UserID:      claims[IdentityKey].(string),
		ESIndexName: claims[IndexKey].(string),
	}
}
func unauthorizedFunc(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
func NewJwtMiddleware(db *database.Database) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         30 * time.Minute,
		MaxRefresh:      30 * time.Minute,
		IdentityKey:     IdentityKey,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandlerFunc,
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			var loginVals database.Login
			if err := ctx.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			//ログイン処理
			result, err := db.Login(&loginVals)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &database.M_User{
				UserID:      result.UserID,
				ESIndexName: result.ESIndexName,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*database.M_User); ok {
				user, err := db.GetUser(v.UserID)
				if err != nil {
					return false
				}
				fmt.Println(user)
				if user.UserID == v.UserID {
					return true
				}
			}
			return false
		},
		Unauthorized: unauthorizedFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
