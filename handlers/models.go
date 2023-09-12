package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zzong12/model-svr-go/model"
)

func models(c *gin.Context) {
	models := model.GetModels()
	c.JSON(200, models)
}
