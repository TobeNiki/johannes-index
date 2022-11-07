package cors

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() gin.HandlerFunc {
	env := os.Getenv("ENV_NAME")
	allowOriginName := ""
	if env == "prd" {
		allowOriginName = "" //今はまだ本番環境がないので、
	} else {
		allowOriginName = "http://localhost:9000" //フロントのローカル環境
	}
	return cors.New(cors.Config{
		AllowOrigins: []string{
			allowOriginName,
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 12 * time.Hour,
	})
}
