package helpers

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

type Merger struct {
	Swagger map[string]interface{}
}

func NewMerger() *Merger {
	merger := new(Merger)
	merger.Swagger = map[string]interface{}{}
	return merger
}

func (m *Merger) AddFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var s1 interface{}
	err = yaml.Unmarshal(content, &s1)
	if err != nil {
		return err
	}

	return m.merge(s1.(map[string]interface{}))
}

func (m *Merger) merge(f map[string]interface{}) error {
	for key, item := range f {
		if i, ok := item.(map[string]interface{}); ok {
			for subKey, subitem := range i {
				if _, ok := m.Swagger[key]; !ok {
					m.Swagger[key] = map[string]interface{}{}
				}

				m.Swagger[key].(map[string]interface{})[subKey] = subitem
			}
		} else {
			m.Swagger[key] = item
		}
	}

	return nil
}

func (m *Merger) Save(fileName string) error {
	res, _ := yaml.Marshal(m.Swagger)

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(res)
	if err != nil {
		return err
	}

	return nil
}
