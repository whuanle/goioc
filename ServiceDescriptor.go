package goioc

import (
	"reflect"
)

// ServiceDescriptor 是注入项的描述，
// 描述如何实例化类型
type ServiceDescriptor struct {
	// 对象名称
	Name string

	// 对象生命周期
	Lifetime ServiceLifetime

	// 对象继承的接口，要注入的接口等
	BaseType reflect.Type

	// 实现对象，实例对象
	ServiceType reflect.Type

	// 已被实例化的对象，存储在内存中
	ServiceInstance interface{}

	// 如何实例化对象，要求返回的必须是对象的指针给接口
	InitHandler func() interface{}
}
