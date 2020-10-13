package controller

import (
	"github.com/abserari/shower/admin/model/mysql"
	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (con *Controller) GetID(ctx *gin.Context) (uint32, error) {
	id, ok := ctx.Get("userID")
	if !ok {
		return 0, errUserIDNotExists
	}

	v, ok := id.(float64)
	if !ok {
		return 0, errUserIDNotValid(id)
	}

	return uint32(v), nil
}

//CheckActive middleware that checks the active
func (con *Controller) CheckActive() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		a, err := con.GetID(ctx)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		active, err := mysql.IsActive(con.db, a)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusConflict, err)
			return
		}

		if !active {
			_ = ctx.AbortWithError(http.StatusLocked, errActive)
			ctx.JSON(http.StatusLocked, gin.H{"status": http.StatusLocked})
			return
		}
	}
}

func (con *Controller) newJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test-pet",
		Key:         []byte("moli-tech-cats-member"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "userID",
		// use data as userID here.
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{
				"userID": data,
			}
		},
		// just get the ID
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["userID"]
		},
		Authenticator: func(ctx *gin.Context) (interface{}, error) {
			return con.Login(ctx)
		},
		// no need to check user valid every time.
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: JWT",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
