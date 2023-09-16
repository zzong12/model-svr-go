package model

import (
	"fmt"
	"log"
	"math"

	"github.com/getumen/go-treelite"
)

var _ Model = (*TreeliteModel)(nil)

type TreeliteModel struct {
	*baseModel
	predictor *treelite.Predictor
}

func (m *TreeliteModel) Load() error {
	predictor, err := treelite.NewPredictor(m.baseModel.Path, 1)
	if err != nil {
		return err
	}
	m.predictor = predictor
	return nil
}

func (m *TreeliteModel) Destory() error {
	return m.predictor.Close()
}

func (m *TreeliteModel) Predict(features [][]float64) ([]float64, error) {
	if m.predictor == nil {
		return nil, fmt.Errorf("predictor is not loaded")
	}
	if len(features) == 0 {
		return nil, fmt.Errorf("features is empty")
	}
	nRow := len(features)
	nCol := len(features[0])
	data := make([]float32, nRow*nCol)
	for i, row := range features {
		for j, col := range row {
			data[i*nCol+j] = float32(col)
		}
	}
	dMatrix, err := treelite.CreateFromMat(data, nRow, nCol, float32(math.NaN()))
	if err != nil {
		log.Fatal(err)
	}
	defer dMatrix.Close()
	result, err := m.predictor.PredictBatch(dMatrix, false, false)
	if err != nil {
		return nil, err
	}
	var ret []float64
	for _, row := range result {
		ret = append(ret, float64(row))
	}
	return ret, nil
}
