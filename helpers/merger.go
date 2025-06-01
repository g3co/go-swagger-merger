package helpers

import (
	"io"
	"os"

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
	if err = yaml.Unmarshal(content, &swaggerMap); err != nil {
		return err
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
	res, err := yaml.Marshal(m.Swagger)
	if err != nil {
		return err
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
