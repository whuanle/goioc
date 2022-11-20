package goioc

import (
	"reflect"
	"sync"
)

var singletonInstanceManager singletonProvider

type singletonDescriptor struct {
	baseType    reflect.Type
	instance    interface{}
	initHandler func(provider IServiceProvider) interface{}
	lock        *sync.Once
}

type singletonProvider struct {
	descriptors map[reflect.Type]*singletonDescriptor
}

// get 获取静态实例
func (descriptor singletonDescriptor) get(provider IServiceProvider) interface{} {
	descriptor.lock.Do(func() {
		if descriptor.instance == nil {
			p := provider.(*ServiceProvider)
			descriptor.instance = p.createObject(descriptor.initHandler(provider), Singleton)
		}
	})
	return descriptor.instance
}

// 注册静态实例
func registerSingletonInstance(baseType reflect.Type, f func(provider IServiceProvider) interface{}) {
	once := &sync.Once{}
	descriptor := singletonDescriptor{
		baseType:    baseType,
		initHandler: f,
		lock:        once,
	}
	singletonInstanceManager.descriptors[baseType] = &descriptor
}

func getSingletonInstance(baseType reflect.Type, provider IServiceProvider) interface{} {
	descriptor := singletonInstanceManager.descriptors[baseType]
	if descriptor == nil {
		return nil
	}
	return descriptor.get(provider)
}
