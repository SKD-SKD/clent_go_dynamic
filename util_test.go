package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testYaml = `apiVersion: v1
kind: Service`

func TestGetKind(t *testing.T) {
	kind, err := GetKind(testYaml)

	assert.Nil(t, err)
	assert.Equal(t, *kind, ServiceKind)
}
