package goioc

import (
	"fmt"
	"reflect"
)

type ServiceProvider struct {
	descriptors map[reflect.Type]ServiceDescriptor
}

// GetService 获取对象实例
func (s *ServiceProvider) GetService(baseType reflect.Type) (*interface{}, error) {
	descriptor, ok := s.descriptors[baseType]
	if !ok {
		return nil, fmt.Errorf("没有找到 %s 对应的实例", baseType.Name())
	}
	if descriptor.Lifetime == Transient {
		// 实例化当前类型
		obj := descriptor.InitHandler()
		// 创建对象并且检查当前结构体是否还有需要被注入的字段
		obj = s.createObject(obj)
		return &obj, nil
	}
	// descriptor.Lifetime == Scope || descriptor.Lifetime == Singleton
	if descriptor.ServiceInstance == nil {
		descriptor.ServiceInstance = s.createObject(descriptor.InitHandler())
	}

	return &descriptor.ServiceInstance, nil
}

// createObject 结构体字段自动注入，
// 递归给需要依赖注入的结构体字段注入实例。
// obj 对应的结构体需要是结构体指针，
// 创建对象后必须返回结构体指针；
func (s *ServiceProvider) createObject(obj interface{}) interface{} {
	sourceType := reflect.TypeOf(obj).Elem()
	if sourceType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("注入的 %t 类型不是结构体！", sourceType))
	}

	// 如果是 interface = Type{} 则不需要处理，否则解开 interface = &Type{} ，获取 Type；
	//必须使用 &obj 而不是 obj，否则无法通过方式为其赋值， reflect.TypeOf(&obj).Elem()
	v := reflect.ValueOf(obj).Elem() // 获取执行的对象
	t := v.Type()

	// 找到需要被依赖注入的字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("injection")
		if tag == "" {
			continue
		}
		if tag == "true" {
			value, err := s.GetService(t.Field(i).Type)
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
