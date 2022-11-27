package services

import (
	"fmt"
	"github.com/whuanle/goioc"
	"reflect"
	"sync"
)

// ServiceCollection 即 IServiceCollection 的实现
type ServiceCollection struct {
	// 服务描述
	descriptors map[reflect.Type]goioc.ServiceDescriptor
	// single对象描述
	singletonDescriptors map[reflect.Type]*SingletonDescriptor
	// 当前在容器中类型数量
	Count int
}

// 基础注入方法
func (s *ServiceCollection) addAny(
	lifetime goioc.ServiceLifetime,
	baseType reflect.Type,
	serviceType reflect.Type,
	f func(provider goioc.IServiceProvider) interface{}) {
	descriptor := goioc.ServiceDescriptor{
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
func getInitHandler(t reflect.Type) func(provider goioc.IServiceProvider) interface{} {
	f := func(provider goioc.IServiceProvider) interface{} {
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

func (s *ServiceCollection) AddService(lifetime goioc.ServiceLifetime, t reflect.Type) {
	checkStructType(t)
	f := getInitHandler(t)
	s.addAny(lifetime, t, t, f)
}

func (s *ServiceCollection) AddServiceHandler(lifetime goioc.ServiceLifetime, t reflect.Type, f func(provider goioc.IServiceProvider) interface{}) {
	checkBaseType(t)
	s.addAny(lifetime, t, t, f)
}

func (s *ServiceCollection) AddServiceOf(lifetime goioc.ServiceLifetime, baseType reflect.Type, serviceType reflect.Type) {
	checkBaseType(baseType)
	checkStructType(serviceType)
	f := getInitHandler(serviceType)
	s.addAny(lifetime, baseType, serviceType, f)
}

func (s *ServiceCollection) AddServiceHandlerOf(
	lifetime goioc.ServiceLifetime,
	baseType reflect.Type,
	serviceType reflect.Type,
	f func(provider goioc.IServiceProvider) interface{}) {
	checkBaseType(baseType)
	checkBaseType(serviceType)
	s.addAny(lifetime, baseType, serviceType, f)
}

// 私有方法

// 添加一个 ServiceDescriptor
func (s *ServiceCollection) add(serviceDescriptor goioc.ServiceDescriptor) {
	if s.descriptors == nil {
		s.descriptors = make(map[reflect.Type]goioc.ServiceDescriptor)
	}
	s.descriptors[serviceDescriptor.BaseType] = serviceDescriptor
	s.Count = len(s.descriptors)
}

func (s *ServiceCollection) get(baseType reflect.Type) *goioc.ServiceDescriptor {
	sd, ok := s.descriptors[baseType]
	if !ok {
		panic(fmt.Sprintf("Type [ %t ] not found", baseType))
	}
	return &sd
}

// 移除一个 ServiceDescriptor
func (s *ServiceCollection) remove(descriptor goioc.ServiceDescriptor) {
	delete(s.descriptors, descriptor.BaseType)
	s.Count = len(s.descriptors)
}

func (s *ServiceCollection) Build() goioc.IServiceProvider {
	// 第一次使用时，初始化单例管理器
	if s.singletonDescriptors == nil {
		s.singletonDescriptors = map[reflect.Type]*SingletonDescriptor{}
	}

	descriptors := make(map[reflect.Type]goioc.ServiceDescriptor)
	onces := make(map[reflect.Type]*sync.Once)

	// 复制集合中的 ServiceDescriptor 到新的容器中，检查
	for i, descriptor := range s.descriptors {
		// 单例模式会被放置到全局实例管理器
		if descriptor.Lifetime == goioc.Singleton {
			s.registerSingletonInstance(descriptor.BaseType, descriptor.InitHandler)
		}
		descriptors[i] = descriptor
		onces[i] = &sync.Once{}
	}

	var services goioc.IServiceProvider
	services = &ServiceProvider{
		descriptors:       descriptors,
		onces:             onces,
		serviceCollection: s,
	}
	return services
}

func (s *ServiceCollection) CopyTo() goioc.IServiceCollection {
	descriptors := make(map[reflect.Type]goioc.ServiceDescriptor)

	for i, descriptor := range s.descriptors {
		descriptors[i] = descriptor
	}
	return &ServiceCollection{
		descriptors: descriptors,
		Count:       len(descriptors),
	}
}

// 静态对象处理
// 注册静态实例
func (s *ServiceCollection) registerSingletonInstance(baseType reflect.Type, f func(provider goioc.IServiceProvider) interface{}) {
	if s.singletonDescriptors[baseType] != nil {
		return
	}
	once := &sync.Once{}
	descriptor := SingletonDescriptor{
		baseType:    baseType,
		initHandler: f,
		lock:        once,
	}
	s.singletonDescriptors[baseType] = &descriptor
}

func (s *ServiceCollection) getSingletonInstance(baseType reflect.Type, provider *ServiceProvider) interface{} {
	descriptor := s.singletonDescriptors[baseType]
	if descriptor == nil {
		return nil
	}
	obj := descriptor.initAndGet(provider)
	return obj
}
