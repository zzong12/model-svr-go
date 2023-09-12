package model

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	modelRepo     = map[string]Model{}
	rwLock        = sync.RWMutex{}
	modelRepoPath string
)

const (
	ENV_MODEL_PATH = "MODEL_REPO"
)

func init() {
	modelRepoPath = os.Getenv(ENV_MODEL_PATH)
	if len(modelRepoPath) == 0 {
		panic(fmt.Sprintf("Can't found environment, %s", ENV_MODEL_PATH))
	}
	LoadModelRepo(modelRepoPath, true)
}

/*
*
model_repo

	-- model_name_1
		-- version_name_1
			-- xxx.so
		-- version_name_2
			-- xxx.so
	-- model_name_2
	...
*/
func LoadModelRepo(path string, autoReload bool) error {
	modelDirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, modelDir := range modelDirs {
		if !modelDir.IsDir() {
			continue
		}
		versionDris, err := os.ReadDir(filepath.Join(path, modelDir.Name()))
		if err != nil {
			log.Println("failed to read dir:", modelDir.Name())
			continue
		}
		for _, versionDir := range versionDris {
			if !versionDir.IsDir() {
				continue
			}
			modelFiles, err := os.ReadDir(filepath.Join(path, modelDir.Name(), versionDir.Name()))
			if err != nil {
				log.Println("failed to read dir:", versionDir.Name())
				continue
			}
			for _, modelFile := range modelFiles {
				if modelFile.IsDir() {
					continue
				}
				modelPath := filepath.Join(path, modelDir.Name(), versionDir.Name(), modelFile.Name())
				modelName := modelDir.Name()
				versionName := versionDir.Name()
				modelSuffix := modelName[len(modelName)-3:]
				modelMd5, _ := getFileMd5(modelPath)

				oldModel := GetModel(modelName, versionName)
				var replaceOld bool
				if oldModel != nil && oldModel.FileMd5() == modelMd5 {
					continue
				} else if oldModel != nil && oldModel.FileMd5() != modelMd5 {
					replaceOld = true
					log.Println("model is updated:", modelPath)
				}

				var model Model
				baseModel := &baseModel{
					Name:    modelName,
					Version: versionName,
					Path:    modelPath,
					Uptime:  time.Now(),
					Md5:     modelMd5,
				}
				if modelSuffix != ".so" {
					model = &TreeliteModel{
						baseModel: baseModel,
					}
				}

				err := model.Load()
				if err != nil {
					log.Println("failed to load model:", modelPath)
					continue
				}

				rwLock.Lock()

				modelRepo[getRepoKey(modelName, versionName)] = model
				rwLock.Unlock()
				log.Printf("load model: %s", modelPath)

				if replaceOld {
					oldModel.Destory()
					log.Println("destory old model:", modelPath)
				}
			}
		}
		if autoReload {
			go func() {
				ticker := time.NewTicker(10 * time.Minute)
				for {
					select {
					case <-ticker.C:
						log.Println("reload model repo:", path)
						LoadModelRepo(path, false)
					}
				}
			}()
		}
	}

	return nil

}

func getFileMd5(filePath string) (string, error) {
	md5Ctx := md5.New()
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(md5Ctx, file)
	if err != nil {
		return "", err
	}
	return string(md5Ctx.Sum(nil)), nil
}

func getRepoKey(modelName string, modelVersion string) string {
	return modelName + "." + modelVersion
}

func GetModel(modelName string, modelVersion string) Model {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if model, ok := modelRepo[getRepoKey(modelName, modelVersion)]; ok {
		return model
	}
	return nil
}

func GetModels() []string {
	rwLock.RLock()
	defer rwLock.RUnlock()
	var models []string
	for modelFullName := range modelRepo {
		models = append(models, modelFullName)
	}
	return nil
}

func Predict(name string, version string, features [][]float64) ([]float64, error) {
	model := GetModel(name, version)
	if model == nil {
		return nil, nil
	}
	return model.Predict(features)
}
