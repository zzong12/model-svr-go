package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zzong12/model-svr-go/handlers"
)

func main() {
	router := gin.Default()
	handlers.Inithandlers(router)

	router.Run(":8000")
}
