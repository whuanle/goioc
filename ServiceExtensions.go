package goioc

import (
	"reflect"
)

// Get 获取对象
func Get[T any](provider IServiceProvider) interface{} {
	t := reflect.TypeOf((*T)(nil)).Elem()

	if t.Kind() == reflect.Interface {
		obj, err := provider.GetService(t)
		if err != nil {
			panic(err)
		}
		return obj
	} else if t.Kind() == reflect.Struct {
		obj, err := provider.GetService(t)
		if err != nil {
			panic(err)
		}
		return obj
	}
	return nil
}

// GetI 根据接口获取对象
func GetI[T interface{}](provider IServiceProvider) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Interface {
		var v T
		return v
	}
	return getInterface[T](provider, t)
}

// GetS 根据结构体获取对象
func GetS[T interface{} | struct{}](provider IServiceProvider) *T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Struct {
		return nil
	}
	return getStruct[T](provider, t)
}

// 接口
func getInterface[T any](provider IServiceProvider, it reflect.Type) T {
	obj, err := provider.GetService(it)
	if err != nil {
		panic(err)
	}

	// 转换为接口
	v := (*obj).(T)
	return v
}

// 结构体
func getStruct[T any](provider IServiceProvider, it reflect.Type) *T {
	obj, err := provider.GetService(it)
	if err != nil {
		panic(err)
	}

	// 转换为结构体指针
	v := (*obj).(*T)
	return v
}
