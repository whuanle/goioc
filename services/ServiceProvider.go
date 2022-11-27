package services

import (
	"fmt"
	"github.com/whuanle/goioc"
	"reflect"
	"sync"
)

type ServiceProvider struct {
	descriptors       map[reflect.Type]goioc.ServiceDescriptor
	onces             map[reflect.Type]*sync.Once
	serviceCollection *ServiceCollection
}

// Dispose 释放所有对象
func (s *ServiceProvider) Dispose() {
	for i, _ := range s.descriptors {
		instance := s.descriptors[i]
		if instance.ScopeInstance != nil {
			if obj, ok := instance.ScopeInstance.(goioc.IDispose); ok {
				obj.Dispose()
			}
			instance.ScopeInstance = nil
		}

		s.descriptors[i] = instance
	}
}

// GetService 获取对象实例
func (s *ServiceProvider) GetService(baseType reflect.Type) (*interface{}, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("error instantiating the object: [ v% ]", err))
		}
	}()
	descriptor, ok := s.descriptors[baseType]
	if !ok {
		return nil, fmt.Errorf("type [ %t ] not found", baseType)
	}
	if descriptor.Lifetime == goioc.Transient {
		obj := descriptor.InitHandler(s)
		// 创建对象并且检查当前结构体是否还有需要被注入的字段
		obj = createObject(s, obj, goioc.Transient)
		return &obj, nil
	}

	// descriptor.Lifetime == Scope
	if descriptor.Lifetime == goioc.Scope {
		if descriptor.ScopeInstance == nil {
			once := s.onces[descriptor.BaseType]
			once.Do(func() {
				obj := descriptor.InitHandler(s)
				descriptor.ScopeInstance = createObject(s, obj, goioc.Scope)
				s.descriptors[baseType] = descriptor
			})
		}
		return &descriptor.ScopeInstance, nil
	}

	// 如果是单例模式，则要找到原始的 collection ，实例化，每次都从 ServiceCollection 中取对象
	if descriptor.Lifetime == goioc.Singleton {
		instance := s.serviceCollection.getSingletonInstance(baseType, s)
		if instance == nil {
			return nil, fmt.Errorf("type [ %t ] not found", baseType)
		}
		return &instance, nil
	}
	panic(fmt.Sprintf("Unrecognized life cycle: [ %v ]", descriptor.Lifetime))
}

// 获取对象，并检测生命周期。
// sourceLifetime：被注入的对象的生命周期
func getService(s *ServiceProvider, baseType reflect.Type, sourceLifetime goioc.ServiceLifetime) (*interface{}, error) {
	descriptor, ok := s.descriptors[baseType]
	if !ok {
		return nil, fmt.Errorf("type [ %t ] not found", baseType)
	}
	if descriptor.Lifetime == goioc.Transient {
		obj := descriptor.InitHandler(s)
		// 创建对象并且检查当前结构体是否还有需要被注入的字段
		obj = createObject(s, obj, descriptor.Lifetime)
		return &obj, nil
	}

	// descriptor.Lifetime == Scope
	if descriptor.Lifetime == goioc.Scope {
		if sourceLifetime == goioc.Singleton {
			panic(fmt.Sprintf("Cannot inject an instance whose lifecycle is scope[ %t ] into singleton", baseType))
		}
		if descriptor.ScopeInstance == nil {
			once := s.onces[descriptor.BaseType]
			once.Do(func() {
				obj := descriptor.InitHandler(s)
				s.descriptors[baseType] = descriptor
				descriptor.ScopeInstance = createObject(s, obj, descriptor.Lifetime)
			})

		}
		return &descriptor.ScopeInstance, nil
	}

	// 如果是单例模式，则要找到原始的 collection ，实例化，每次都从 ServiceCollection 中取对象
	if descriptor.Lifetime == goioc.Singleton {
		instance := s.serviceCollection.getSingletonInstance(baseType, s)
		if instance == nil {
			return nil, fmt.Errorf("type [ %t ] not found", baseType)
		}
		return &instance, nil
	}
	panic(fmt.Sprintf("Unrecognized life cycle: [ %v ]", descriptor.Lifetime))
}

// createObject 结构体字段自动注入，
// 递归给需要依赖注入的结构体字段注入实例。
// obj 对应的结构体需要是结构体指针，
// 创建对象后必须返回结构体指针；
func createObject(s *ServiceProvider, obj interface{}, lifetime goioc.ServiceLifetime) interface{} {
	sourceType := reflect.TypeOf(obj).Elem()
	if sourceType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[ %t ] is not an interface or struct", sourceType))
	}

	// 如果是 interface = Type{} 则不需要处理，否则解开 interface = &Type{} ，获取 Type；
	//必须使用 &obj 而不是 obj，否则无法通过方式为其赋值， reflect.TypeOf(&obj).Elem()
	v := reflect.ValueOf(obj).Elem() // 获取执行的对象
	t := v.Type()

	// 找到需要被依赖注入的字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("ioc")
		if tag == "" {
			continue
		}
		if tag == "true" {
			// 字段类型，如果字段类型是指针，则需要解开指针
			fieldSourceType := field.Type
			if fieldSourceType.Kind() == reflect.Ptr {
				fieldSourceType = fieldSourceType.Elem()
			}
			value, err := getService(s, fieldSourceType, lifetime)
			if err != nil {
				panic(err)
			}

			// 赋值
			rValue := reflect.ValueOf(*value)
			switch field.Type.Kind() {
			case reflect.Interface:
				break
			case reflect.Ptr:
				break
			case reflect.Struct:
				rValue = rValue.Elem()
			}

			v.Field(i).Set(rValue)
		}
	}
	return obj
}
