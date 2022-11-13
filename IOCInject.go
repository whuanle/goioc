package goioc

import "reflect"

// AddScoped 注册对象为域(Scope)生命周期
func AddScoped[T any](con IServiceCollection) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddScoped(t)
}

// AddScopedHandler 注册对象为域(Scope)生命周期，并自定义如何初始化实例
func AddScopedHandler[T any](con IServiceCollection, f func() interface{}) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddScopedHandler(t, f)
}

// AddScopedForm 注册对象为域(Scope)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType
func AddScopedForm[I any, T any](con IServiceCollection) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddScopedForm(i, t)
}

// AddScopedHandlerForm 注册对象为域(Scope)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
func AddScopedHandlerForm[I any, T any](con IServiceCollection, f func() interface{}) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddScopedHandlerForm(i, t, f)
}

// AddSingleton 注册对象为单例(Singleton)
func AddSingleton[T any](con IServiceCollection) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddSingleton(t)
}

// AddSingletonHandler 注册对象为单例(Singleton)，并自定义如何初始化实例
func AddSingletonHandler[T any](con IServiceCollection, f func() interface{}) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddSingletonHandler(t, f)
}

// AddSingletonForm 注注册对象为单例(Singleton)，注册接口或父类型及其实现，serviceType 必须实现了 baseType
func AddSingletonForm[I any, T any](con IServiceCollection) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddSingletonForm(i, t)
}

// AddSingletonHandlerForm 注册对象为单例(Singleton)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
func AddSingletonHandlerForm[I any, T any](con IServiceCollection, f func() interface{}) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddSingletonHandlerForm(i, t, f)
}

// AddTransient 注册对象为短暂(Transient)生命周期
func AddTransient[T any](con IServiceCollection) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddTransient(t)
}

// AddTransientHandler 注册对象为短暂(Transient)，并自定义如何初始化实例
func AddTransientHandler[T any](con IServiceCollection, f func() interface{}) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddTransientHandler(t, f)
}

// AddTransientForm 注注册对象为短暂(Transient)，注册接口或父类型及其实现，serviceType 必须实现了 baseType
func AddTransientForm[I any, T any](con IServiceCollection) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddTransientForm(i, t)
}

// AddTransientHandlerForm 注册对象为短暂(Transient)生命周期，注册接口或父类型及其实现，serviceType 必须实现了 baseType，并自定义如何初始化实例
func AddTransientHandlerForm[I any, T any](con IServiceCollection, f func() interface{}) {
	i := reflect.TypeOf((*I)(nil)).Elem()
	t := reflect.TypeOf((*T)(nil)).Elem()
	con.AddTransientHandlerForm(i, t, f)
}
