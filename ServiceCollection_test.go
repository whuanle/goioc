package goioc

import (
	"reflect"
	"testing"
)

// 普通的测试
func TestCreateCollection(t *testing.T) {
	var sc IServiceCollection = &ServiceCollection{}

	var is IServiceCollection
	sc.AddScopedForm(reflect.TypeOf(is), reflect.TypeOf(ServiceCollection{}))

	_ = sc.Build()
}
