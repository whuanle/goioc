package goioc

import (
	"fmt"
	"reflect"
)

// ServiceCollection 即 IServiceCollection 的实现
type ServiceCollection struct {
	// 服务描述
	descriptors map[reflect.Type]ServiceDescriptor
	// 当前在容器中类型数量
	Count int
}

// 基础注入方法
func (s *ServiceCollection) addAny(
	lifetime ServiceLifetime,
	baseType reflect.Type,
	serviceType reflect.Type,
	f func(provider IServiceProvider) interface{}) {
	descriptor := ServiceDescriptor{
		Name:        serviceType.Name(),
		BaseType:    baseType,
		ServiceType: serviceType,
		Lifetime:    lifetime,
		InitHandler: f,
	}
	s.add(descriptor)
}

// 实现 IServiceCollection 接口

// 默认通过反射实例化对象的方式
func getInitHandler(t reflect.Type) func(provider IServiceProvider) interface{} {
	f := func(provider IServiceProvider) interface{} {
		// reflect 先实例化此类型，获取了当前对象的指针，获取原始对象
		return reflect.New(t).Interface()
	}
	return f
}

// 检查类型是否为接口或结构体
func checkBaseType(t reflect.Type) {
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Struct {
		return
	}
	panic(fmt.Sprintf("[ %t ] is not an interface or struct", t))
}

// 检查是否为结构体
func checkStructType(t reflect.Type) {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[ %t ] is not an struct", t))
	}
}

func (s *ServiceCollection) AddService(lifetime ServiceLifetime, t reflect.Type) {
	checkStructType(t)
	f := getInitHandler(t)
	s.addAny(lifetime, t, t, f)
}

func (s *ServiceCollection) AddServiceHandler(lifetime ServiceLifetime, t reflect.Type, f func(provider IServiceProvider) interface{}) {
	checkBaseType(t)
	s.addAny(lifetime, t, t, f)
}

func (s *ServiceCollection) AddServiceOf(lifetime ServiceLifetime, baseType reflect.Type, serviceType reflect.Type) {
	checkBaseType(baseType)
	checkStructType(serviceType)
	f := getInitHandler(serviceType)
	s.addAny(lifetime, baseType, serviceType, f)
}

func (s *ServiceCollection) AddServiceHandlerOf(
	lifetime ServiceLifetime,
	baseType reflect.Type,
	serviceType reflect.Type,
	f func(provider IServiceProvider) interface{}) {
	checkBaseType(baseType)
	checkBaseType(serviceType)
	s.addAny(lifetime, baseType, serviceType, f)
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
		panic(fmt.Sprintf("Type [ %t ] not found", baseType))
	}
	return &sd
}

// 移除一个 ServiceDescriptor
func (s *ServiceCollection) remove(descriptor ServiceDescriptor) {
	delete(s.descriptors, descriptor.BaseType)
	s.Count = len(s.descriptors)
}

func (s *ServiceCollection) Build() IServiceProvider {
	descriptors := make(map[reflect.Type]ServiceDescriptor)

	// 复制集合中的 ServiceDescriptor 到新的容器中，检查
	for i, descriptor := range s.descriptors {
		// 单例模式会被放置到全局实例管理器
		if descriptor.Lifetime == Singleton {
			registerSingletonInstance(descriptor.BaseType, descriptor.InitHandler)
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
