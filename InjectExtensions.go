package goioc

import "reflect"

// AddService 注册对象
func AddService[T any](con IServiceCollection, lifetime ServiceLifetime) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddService(lifetime, t)
}

// AddServiceHandler 注册对象，并自定义如何初始化实例
func AddServiceHandler[T any](con IServiceCollection, lifetime ServiceLifetime, f func(provider IServiceProvider) interface{}) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddServiceHandler(lifetime, t, f)
}

// AddServiceOf 注册对象，注册接口或父类型及其实现，serviceType 必须实现了 baseType
func AddServiceOf[I any, T any](con IServiceCollection, lifetime ServiceLifetime) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddServiceOf(lifetime, i, t)
}

// AddServiceHandlerOf 注册对象，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
func AddServiceHandlerOf[I any, T any](con IServiceCollection, lifetime ServiceLifetime, f func(provider IServiceProvider) interface{}) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddServiceHandlerOf(lifetime, i, t, f)
}
