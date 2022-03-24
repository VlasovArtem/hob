package dependency

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestObject struct{}

type TestInterface interface {
	hello() string
}

func (t TestObject) hello() string {
	return "test"
}

func Test(t *testing.T) {
	provider := NewDependenciesProvider()

	provider.Add(TestObject{})

	dependency := FindRequiredDependency[TestObject, TestInterface](provider)

	assert.NotNil(t, dependency)
	assert.Equal(t, "test", dependency.hello())
}
