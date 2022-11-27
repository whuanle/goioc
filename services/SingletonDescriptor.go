package services

import (
	"github.com/whuanle/goioc"
	"reflect"
	"sync"
)

// SingletonDescriptor 静态对象描述
type SingletonDescriptor struct {
	baseType    reflect.Type
	instance    interface{}
	initHandler func(provider goioc.IServiceProvider) interface{}
	lock        *sync.Once
}

// 初始化
func (descriptor *SingletonDescriptor) initAndGet(provider *ServiceProvider) interface{} {
	descriptor.lock.Do(func() {
		if descriptor.instance == nil {
			descriptor.instance = createObject(provider, descriptor.initHandler(provider), goioc.Singleton)
		}
	})
	return descriptor.instance
}
