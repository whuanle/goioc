package goioc

import "reflect"

// IServiceProvider 依赖注入提供器，
// 将类型实例化为对象。
type IServiceProvider interface {
	// GetService 获取你需要的服务实例
	GetService(baseType reflect.Type) (*interface{}, error)
	// Dispose 释放当前容器的 Scope 对象
	Dispose()
}
