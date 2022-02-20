package goioc

import "reflect"

// IServiceCollection 是依赖注入对象容器接口，
// 将类型注入到容器中
type IServiceCollection interface {

	// AddScoped 注册对象为域(Scope)生命周期
	AddScoped(t reflect.Type)
	// AddScopedHandler 注册对象为域(Scope)生命周期，并自定义如何初始化实例
	AddScopedHandler(t reflect.Type, f func() interface{})
	// AddScopedForm 注册对象为域(Scope)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType
	AddScopedForm(baseType reflect.Type, serviceType reflect.Type)
	// AddScopedHandlerForm 注册对象为域(Scope)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
	AddScopedHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{})

	// AddSingleton 注册对象为单例(Singleton)
	AddSingleton(t reflect.Type)
	// AddSingletonHandler 注册对象为单例(Singleton)，并自定义如何初始化实例
	AddSingletonHandler(t reflect.Type, f func() interface{})
	// AddSingletonForm 注注册对象为单例(Singleton)，注册接口或父类型及其实现，serviceType 必须实现了 baseType
	AddSingletonForm(baseType reflect.Type, serviceType reflect.Type)
	// AddSingletonHandlerForm 注册对象为单例(Singleton)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
	AddSingletonHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{})

	// AddTransient 注册对象为短暂(Transient)生命周期
	AddTransient(t reflect.Type)
	// AddTransientHandler 注册对象为短暂(Transient)，并自定义如何初始化实例
	AddTransientHandler(t reflect.Type, f func() interface{})
	// AddTransientForm 注注册对象为短暂(Transient)，注册接口或父类型及其实现，serviceType 必须实现了 baseType
	AddTransientForm(baseType reflect.Type, serviceType reflect.Type)
	// AddTransientHandlerForm 注册对象为短暂(Transient)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
	AddTransientHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{})

	// CopyTo 复制当前容器的所有对象，生成新的容器
	CopyTo() IServiceCollection
	// 	Build() 构建依赖注入服务提供器 IServiceProvider
	Build() IServiceProvider
}
