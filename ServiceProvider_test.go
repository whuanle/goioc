package goioc

import (
	"fmt"
	"reflect"
	"testing"
)

// 定义接口和结构体

type IAnimal interface {
	Println(s string)
}
type Dog struct {
}

func (my Dog) Println(s string) {
	fmt.Println(s)
}

type Animal struct {
	Dog IAnimal `injection:"true"`
}

func TestInterface(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()

	// 创建容器
	var collection IServiceCollection = &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	collection.AddScopedForm(imy, my)

	// 构建服务 Provider
	serviceProvider := collection.Build()

	// 获取对象
	// *interface{} = &Dog{}，因此需要处理指针
	obj, err := serviceProvider.GetService(imy)

	if err != nil {
		t.Error(err)
	}

	fmt.Println("obj 类型是", reflect.ValueOf(obj).Type())

	animal := (*obj).(IAnimal)
	// 	a := (*obj).(*Dog)
	animal.Println("测试")
}

func TestStruct_Field(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()
	animalType := reflect.TypeOf((*Animal)(nil)).Elem()

	// 创建容器
	var collection IServiceCollection = &ServiceCollection{}
	collection.AddScopedForm(imy, my)
	collection.AddScoped(animalType)

	// 构建服务 Provider
	serviceProvider := collection.Build()

	// 获取对象
	// *interface{}
	obj, err := serviceProvider.GetService(animalType)
	if err != nil {
		t.Error(err)
	}

	// *interface{} -> Animal
	fmt.Println(*obj)
	animal := (*obj).(*Animal)
	animal.Dog.Println("测试2")
}

func TestScopeLifetime(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()

	// 创建容器
	var collection IServiceCollection = &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	collection.AddScopedForm(imy, my)

	// 构建服务 Provider
	serviceProvider := collection.Build()

	// 获取对象
	// *interface{} = &Dog{}，因此需要处理指针
	obj1, _ := serviceProvider.GetService(imy)
	obj2, _ := serviceProvider.GetService(imy)

	fmt.Printf("obj1 = %p,obj2 = %p", (*obj1).(*Dog), (*obj2).(*Dog))
	if fmt.Sprintf("%p",(*obj1).(*Dog)) != fmt.Sprintf("%p",(*obj2).(*Dog)){
		t.Error("两个对象不是同一个")
	}
}
