package handlers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Inithandlers(engine *gin.Engine) {
	pprof.Register(engine)
	engine.POST("/predict", predict)
}
