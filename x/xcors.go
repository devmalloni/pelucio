package x

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetCorsConfig(allowOrigins []string) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}

	// corsConfig.AllowOriginFunc = func(origin string) bool {
	// 	match, _ := regexp.MatchString(`http(s)?://localhost:[0-9]+`, origin)
	// 	return match
	// }

	return cors.New(corsConfig)
}
