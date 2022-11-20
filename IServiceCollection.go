package goioc

import "reflect"

// IServiceCollection 是依赖注入对象容器接口，
// 将类型注入到容器中
type IServiceCollection interface {
	// AddService 注册一个服务，t 必须是结构体类型，由于容器对其进行实例化
	AddService(lifetime ServiceLifetime, t reflect.Type)
	// AddServiceHandler 注册一个服务，t 可以是接口或结构体，由开发者决定如何返回实例
	AddServiceHandler(lifetime ServiceLifetime, t reflect.Type, f func(provider IServiceProvider) interface{})
	// AddServiceOf 注册一个服务，baseType 可以是接口或结构体，serviceType 必须是结构体
	AddServiceOf(lifetime ServiceLifetime, baseType reflect.Type, serviceType reflect.Type)
	// AddServiceHandlerOf 注册一个服务，serviceType 必须继承了 baseType，由开发者决定如何返回实例
	AddServiceHandlerOf(lifetime ServiceLifetime, baseType reflect.Type, serviceType reflect.Type, f func(provider IServiceProvider) interface{})

	// CopyTo 复制当前容器的所有注入信息，生成新的容器
	CopyTo() IServiceCollection
	// 	Build() 构建依赖注入服务提供器 IServiceProvider
	Build() IServiceProvider
}
