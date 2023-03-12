package gconfig

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadYamlCfg(path string, v interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("LoadYamlCfg: read file:%v err:%v", path, err))
	}

	var m = make(map[string]interface{})
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(fmt.Errorf("LoadYamlCfg: unmarshal err:%v", err))
	}

	rv := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr {
		panic("LoadYamlCfg: type v must be Ptr")
	}

	checkYaml([]string{rv.Elem().String()}, rv.Elem(), m)

	err = yaml.Unmarshal(data, v)
	if err != nil {
		panic(fmt.Errorf("LoadYamlCfg: unmarshal err:%v", err))
	}
}

func checkYaml(root []string, vt reflect.Type, m map[string]interface{}) {
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		fieldName := field.Name
		if name := field.Tag.Get("yaml"); len(name) > 0 {
			fieldName = name
		}

		cpRoot := make([]string, len(root))
		copy(cpRoot, root)
		cpRoot = append(cpRoot, fieldName)
		if v, ok := m[fieldName]; !ok || v == nil {
			panic(fmt.Errorf("LoadYamlCfg: field %v not set", strings.Join(cpRoot, ".")))
		}

		if field.Type.Kind() == reflect.Struct {
			if mv, ok := m[fieldName].(map[string]interface{}); ok {
				checkYaml(cpRoot, field.Type, mv)

			} else {
				panic(fmt.Errorf("LoadYamlCfg: field %v is struct", strings.Join(cpRoot, ".")))
			}
		}
	}
}
