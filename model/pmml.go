package model

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/asafschers/goscore"
)

// TODO: implement

var _ Model = (*PmmlModel)(nil)

type PmmlModel struct {
	*baseModel
	predicate *goscore.GradientBoostedModel
}

func (m *PmmlModel) Load() error {
	treeXml, err := os.ReadFile(m.Path)
	if err != nil {
		return fmt.Errorf("pmml is not loaded")
	}
	var model goscore.GradientBoostedModel
	err = xml.Unmarshal([]byte(treeXml), &model)
	if err != nil {
		return err
	}
	m.predicate = &model
	return nil
}

func (m *PmmlModel) Destory() error {
	return nil
}

func (m *PmmlModel) Predict(features [][]float64) ([]float64, error) {
	// ret := make([]float64, len(features))
	// for i, feature := range features {
	// 	data := make(map[string]interface{})
	// 	for i, attr := range m.predicate.Score() {
	// 		data[attr.Name.Local] = feature[i]
	// 	}
	// 	score, err := m.predicate.TraverseTree(data)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ret[i] = score
	// }
	// return ret, nil
}
