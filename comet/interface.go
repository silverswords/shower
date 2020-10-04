package comet

import (
	"github.com/gin-gonic/gin"
)

type comet interface {
	Init(Metadata)
	RegisterRouter(gin.IRouter)
}

type Metadata map[string]interface{}
