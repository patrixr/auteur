package common

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type Metadata map[string]interface{}

func MetaToStruct[T any](m Metadata) (T, error) {
	var result T
	err := mapstructure.Decode(m, &result)
	return result, err
}

func (m *Metadata) PrettyPrint() string {
	yamlData, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(yamlData)
}
