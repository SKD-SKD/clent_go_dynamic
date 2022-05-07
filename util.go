package kubernetes

import (
	"code.cargurus.com/platform/glados/pkg/util"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func HasNamespace(name string, namespaces *v1.NamespaceList) bool {
	for _, namespace := range namespaces.Items {
		if namespace.Name == name {
			return true
		}
	}
	return false
}

func IsRecognizedKind(kind Kind) bool {
	if _, exists := Kinds[kind]; exists {
		return true
	}
	return false
}

func GetKind(definition YAML) (*Kind, error) {
	var d map[string]interface{}
	err := yaml.Unmarshal([]byte(definition), &d)
	if err != nil {
		return nil, err
	}
	if kind, ok := d["kind"].(string); ok {
		if IsRecognizedKind(kind) {
			return &kind, nil
		}
	}
	return nil, util.ErrNotFound
}
