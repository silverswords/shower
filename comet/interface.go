package comet

import (
	"github.com/gin-gonic/gin"
)

type Comet interface {
	RegisterRouter(gin.IRouter)
}
