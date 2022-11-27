package services

import (
	"github.com/whuanle/goioc"
	"testing"
)

// 普通的测试
func TestCreateCollection(t *testing.T) {
	var sc goioc.IServiceCollection = &ServiceCollection{}

	// var is goioc.IServiceCollection
	// sc.AddServiceOf(goioc.Scope, reflect.TypeOf(is), reflect.TypeOf(ServiceCollection{}))
	goioc.AddServiceOf[goioc.IServiceCollection, ServiceCollection](sc, goioc.Scope)
	p := sc.Build()
	obj := goioc.Get[goioc.IServiceCollection](p)
	if obj == nil {
		t.Errorf("service is nil!")
	}
}
