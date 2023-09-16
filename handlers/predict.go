package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zzong12/model-svr-go/model"
)

type (
	PredictRequest struct {
		ModelName    string        `json:"model_name"`
		ModelVersion string        `json:"model_version"`
		RequestItems []RequestItem `json:"items"`
	}
	RequestItem struct {
		Id       string    `json:"id"`
		Features []float64 `json:"features"`
	}

	PredictResponse struct {
		Items []ResponseItem `json:"items"`
	}
	ResponseItem struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
)

func predict(c *gin.Context) {
	var req PredictRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	features := make([][]float64, 0, len(req.RequestItems))
	for _, item := range req.RequestItems {
		features = append(features, item.Features)
	}
	res, err := model.Predict(req.ModelName, req.ModelVersion, features)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	respItems := make([]ResponseItem, 0, len(req.RequestItems))
	for i, item := range req.RequestItems {
		respItems = append(respItems, ResponseItem{
			Id:    item.Id,
			Score: res[i],
		})
	}
	c.JSON(http.StatusOK, PredictResponse{
		Items: respItems,
	})

}
