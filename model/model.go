package model

import "time"

type Model interface {
	FileMd5() string
	Load() error
	Destory() error
	Predict(features [][]float64) ([]float64, error)
}

type baseModel struct {
	Path    string
	Name    string
	Version string
	Md5     string
	Uptime  time.Time
}

func (m *baseModel) FileMd5() string {
	return m.Md5
}
