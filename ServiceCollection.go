package goioc

import (
	"fmt"
	"reflect"
)

// ServiceCollection 即 IServiceCollection 的实现
type ServiceCollection struct {
	descriptors map[reflect.Type]ServiceDescriptor
	// 当前在容器中类型数量
	Count int
}

// 基础注入方法
func (s *ServiceCollection) addAny(baseType reflect.Type, serviceType reflect.Type, lifetime ServiceLifetime, f func() interface{}) {
	// 要去重
	descriptor := ServiceDescriptor{
		Name:        serviceType.Name(),
		BaseType:    baseType,
		ServiceType: serviceType,
		Lifetime:    lifetime,
		InitHandler: f,
	}
	s.add(descriptor)
}

// 下面是实现 IServiceCollection 接口

// 获取实例化此类型的匿名函数
func getInitHandler(t reflect.Type) func() interface{} {
	f := func() interface{} {
		// reflect 先实例化此类型，获取了当前对象的指针，获取原始对象
		return reflect.New(t).Interface()
	}
	return f
}

func checkBaseType(t reflect.Type) {
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Struct {
		return
	}
	panic(fmt.Sprintf("%t 既不是接口也不是结构体，无法注入到容器中", t))
}

func checkServiceType(t reflect.Type) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%t 不是结构体，无法注入到容器中", t))
	}

}

func checkType(baseType reflect.Type, serviceType reflect.Type) {
	checkBaseType(baseType)
	checkServiceType(serviceType)
	if !serviceType.Implements(baseType) {
		panic(fmt.Sprintf("%t 不实现 %t", serviceType, baseType))
	}
}

func (s *ServiceCollection) AddScoped(t reflect.Type) {
	checkServiceType(t)
	f := getInitHandler(t)
	s.addAny(t, t, Scope, f)
}

func (s *ServiceCollection) AddScopedHandler(t reflect.Type, f func() interface{}) {
	checkServiceType(t)
	s.addAny(t, t, Scope, f)
}

func (s *ServiceCollection) AddScopedForm(baseType reflect.Type, serviceType reflect.Type) {
	checkType(baseType, serviceType)
	f := getInitHandler(serviceType)
	s.addAny(baseType, serviceType, Scope, f)
}

func (s *ServiceCollection) AddScopedHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{}) {
	checkType(baseType, serviceType)
	s.addAny(baseType, serviceType, Scope, f)
}

func (s *ServiceCollection) AddSingleton(t reflect.Type) {
	checkServiceType(t)
	f := getInitHandler(t)
	s.addAny(t, t, Singleton, f)
}

func (s *ServiceCollection) AddSingletonHandler(t reflect.Type, f func() interface{}) {
	checkServiceType(t)
	s.addAny(t, t, Singleton, f)
}

func (s *ServiceCollection) AddSingletonForm(baseType reflect.Type, serviceType reflect.Type) {
	checkType(baseType, serviceType)
	f := getInitHandler(serviceType)
	s.addAny(baseType, serviceType, Singleton, f)
}

func (s *ServiceCollection) AddSingletonHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{}) {
	checkType(baseType, serviceType)
	s.addAny(baseType, serviceType, Singleton, f)
}

func (s *ServiceCollection) AddTransient(t reflect.Type) {
	checkServiceType(t)
	f := getInitHandler(t)
	s.addAny(t, t, Transient, f)
}

func (s *ServiceCollection) AddTransientHandler(t reflect.Type, f func() interface{}) {
	checkServiceType(t)
	s.addAny(t, t, Transient, f)
}

func (s *ServiceCollection) AddTransientForm(baseType reflect.Type, serviceType reflect.Type) {
	checkType(baseType, serviceType)
	f := getInitHandler(serviceType)
	s.addAny(baseType, serviceType, Transient, f)
}

func (s *ServiceCollection) AddTransientHandlerForm(baseType reflect.Type, serviceType reflect.Type, f func() interface{}) {
	checkType(baseType, serviceType)
	s.addAny(baseType, serviceType, Transient, f)
}

// 私有方法

// 添加一个 ServiceDescriptor
func (s *ServiceCollection) add(serviceDescriptor ServiceDescriptor) {
	if s.descriptors == nil {
		s.descriptors = make(map[reflect.Type]ServiceDescriptor)
	}
	s.descriptors[serviceDescriptor.BaseType] = serviceDescriptor
	s.Count = len(s.descriptors)
}

func (s *ServiceCollection) get(baseType reflect.Type) *ServiceDescriptor {
	sd, ok := s.descriptors[baseType]
	if !ok {
		panic(fmt.Sprintf("未找到 %v 实例", baseType))
	}
	return &sd
}

// 移除一个 ServiceDescriptor
func (s *ServiceCollection) remove(descriptor ServiceDescriptor) {
	delete(s.descriptors, descriptor.BaseType)
	s.Count = len(s.descriptors)
}

func (s *ServiceCollection) Build() IServiceProvider {
	// scoped
	descriptors := make(map[reflect.Type]ServiceDescriptor)

	// 复制集合中的所有对象到新的容器中，并且对每个 Scope 的对象实例化
	for i, descriptor := range s.descriptors {
		// 单例模式提前实例化，也就是常驻内存
		if descriptor.Lifetime == Singleton {
			if descriptor.ServiceInstance == nil {
				descriptor.ServiceInstance = descriptor.InitHandler() // InitHandler 必须传递了对象指针给接口
			}
		}
		descriptors[i] = descriptor
	}

	var services IServiceProvider
	services = &ServiceProvider{
		descriptors:       descriptors,
		serviceCollection: s,
	}
	return services
}

func (s *ServiceCollection) CopyTo() IServiceCollection {
	descriptors := make(map[reflect.Type]ServiceDescriptor)

	for i, descriptor := range s.descriptors {
		descriptors[i] = descriptor
	}
	return &ServiceCollection{
		descriptors: descriptors,
		Count:       len(descriptors),
	}
}
