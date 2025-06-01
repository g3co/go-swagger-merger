package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type Merger struct {
	Swagger map[string]any
}

func NewMerger() *Merger {
	merger := new(Merger)
	merger.Swagger = map[string]any{}
	return merger
}

func (m *Merger) AddFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var swaggerMap any

	fileExt := filepath.Ext(file)
	switch fileExt {
	case ".yaml", ".yml":
		if err = yaml.Unmarshal(content, &swaggerMap); err != nil {
			return err
		}
	case ".json":
		if err = json.Unmarshal(content, &swaggerMap); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", fileExt)
	}

	merge(m.Swagger, swaggerMap.(map[string]any))

	return nil
}

func merge(a, b map[string]any) {
	if a == nil {
		return
	}

	for key, item := range b {
		if i, ok := item.(map[string]any); ok {
			if _, ok := a[key]; ok {
				merge(a[key].(map[string]any), i)
			} else {
				a[key] = i
			}
		} else {
			a[key] = item
		}
	}
}

func (m *Merger) Save(fileName string) error {
	res := []byte{}
	var err error

	fileExt := filepath.Ext(fileName)
	switch fileExt {
	case ".yaml", ".yml":
		if res, err = yaml.Marshal(m.Swagger); err != nil {
			return err
		}
	case ".json":
		if res, err = json.Marshal(m.Swagger); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", fileExt)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(res); err != nil {
		return err
	}

	return nil
}
