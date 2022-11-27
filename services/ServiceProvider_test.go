package services

import (
	"fmt"
	"github.com/whuanle/goioc"
	"reflect"
	"testing"
)

// 定义接口和结构体

type IAnimal interface {
	Println(s string)
}

type Dog struct {
	Id int
}

func (my Dog) Println(s string) {
	fmt.Println(s)
}

type Animal struct {
	Dog IAnimal `ioc:"true"`
}

func TestServiceCollection_AddService(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddService[Dog](sc, goioc.Scope)
	p := sc.Build()
	obj := goioc.GetS[Dog](p)
	obj.Println("test")
}

func TestServiceCollection_AddServiceHandler(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddServiceHandler[Dog](sc, goioc.Scope, func(provider goioc.IServiceProvider) interface{} {
		return &Dog{
			Id: 1,
		}
	})
	p := sc.Build()
	obj := goioc.GetS[Dog](p)
	if obj == nil {
		t.Errorf("service is nil!")
	}
	obj.Println("test")
}

func TestServiceCollection_AddServiceOf(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddServiceOf[IAnimal, Dog](sc, goioc.Scope)
	p := sc.Build()
	animal1 := goioc.GetI[IAnimal](p)
	if animal1 == nil {
		t.Errorf("service is nil!")
	}
	animal1.Println("test")

	animal2 := goioc.GetI[IAnimal](p)
	if animal2 == nil {
		t.Errorf("service is nil!")
	}

	if animal1 != animal2 {
		t.Errorf("animal1 != animal2")
	}
}

func TestService_Lifetime(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddServiceHandler[Dog](sc, goioc.Singleton, func(provider goioc.IServiceProvider) interface{} {
		return &Dog{
			Id: 2,
		}
	})

	goioc.AddServiceHandlerOf[IAnimal, Dog](sc, goioc.Scope, func(provider goioc.IServiceProvider) interface{} {
		return &Dog{
			Id: 2,
		}
	})

	p := sc.Build()

	a := goioc.GetI[IAnimal](p)

	if v := a.(*Dog); v == nil {
		t.Errorf("service is nil!")
	}
	v := a.(*Dog)
	if v.Id != 2 {
		t.Errorf("Life cycle error")
	}
	v.Id = 3
	// 重复获取的必定是同一个对象
	aa := goioc.GetI[IAnimal](p)
	v = aa.(*Dog)
	if v.Id != 3 {
		t.Errorf("Life cycle error")
	}

	// 重新构建的，scope 不是同一个对象
	pp := sc.Build()
	aaa := goioc.GetI[IAnimal](pp)
	v = aaa.(*Dog)
	if v.Id != 2 {
		t.Errorf("Life cycle error")
	}

	b := goioc.GetS[Dog](p)
	if b.Id != 2 {
		t.Errorf("Life cycle error")
	}

	b.Id = 3

	bb := goioc.GetS[Dog](p)
	if b.Id != bb.Id {
		t.Errorf("Life cycle error")
	}
	ppp := sc.Build()

	bbb := goioc.GetS[Dog](ppp)
	if b.Id != bbb.Id {
		t.Errorf("Life cycle error")
	}
}

func TestServiceProvider_GetService(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddService[Dog](sc, goioc.Scope)
	goioc.AddServiceOf[IAnimal, Dog](sc, goioc.Scope)
	p := sc.Build()

	a := goioc.Get[IAnimal](p)
	b := goioc.Get[Dog](p)
	c := goioc.GetI[IAnimal](p)
	d := goioc.GetS[Dog](p)

	if a == nil {
		t.Errorf("service is nil!")
	}
	if b == nil {
		t.Errorf("service is nil!")
	}
	if c == nil {
		t.Errorf("service is nil!")
	}
	if d == nil {
		t.Errorf("service is nil!")
	}
}

func TestInjectField(t *testing.T) {
	sc := &ServiceCollection{}
	goioc.AddServiceHandlerOf[IAnimal, Dog](sc, goioc.Scope, func(provider goioc.IServiceProvider) interface{} {
		return &Dog{
			Id: 666,
		}
	})
	goioc.AddService[Animal](sc, goioc.Scope)

	p := sc.Build()
	a := goioc.GetS[Animal](p)
	if dog := a.Dog.(*Dog); dog.Id != 666 {
		t.Errorf("service is nil!")
	}
}

func TestReflect_Inject(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()

	// 创建容器
	sc := &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	sc.AddServiceOf(goioc.Scope, imy, my)

	// 构建服务 Provider
	serviceProvider := sc.Build()

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

func TestReflect_Struct_Field(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()
	animalType := reflect.TypeOf((*Animal)(nil)).Elem()

	// 创建容器
	sc := &ServiceCollection{}
	sc.AddServiceOf(goioc.Scope, imy, my)
	sc.AddService(goioc.Scope, animalType)

	// 构建服务 Provider
	serviceProvider := sc.Build()

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
	var sc goioc.IServiceCollection = &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	sc.AddServiceOf(goioc.Scope, imy, my)

	// 构建服务 Provider
	serviceProvider := sc.Build()

	// 获取对象
	// *interface{} = &Dog{}，因此需要处理指针
	obj1, _ := serviceProvider.GetService(imy)
	obj2, _ := serviceProvider.GetService(imy)

	fmt.Printf("obj1 = %p,obj2 = %p\r\n", (*obj1).(*Dog), (*obj2).(*Dog))
	if fmt.Sprintf("%p", (*obj1).(*Dog)) != fmt.Sprintf("%p", (*obj2).(*Dog)) {
		t.Error("两个对象不是同一个")
	}

	// Singleton
	t1 := reflect.TypeOf((*Animal1)(nil)).Elem()
	sc.AddService(goioc.Singleton, t1)

	serviceProvider1 := sc.Build()
	v1, _ := serviceProvider1.GetService(t1)

	serviceProvider2 := sc.Build()
	v2, _ := serviceProvider2.GetService(t1)

	fmt.Printf("v1 = %p,v2 = %p\r\n", (*v1).(*Animal1), (*v2).(*Animal1))
	if fmt.Sprintf("%p", (*v1).(*Animal1)) != fmt.Sprintf("%p", (*v2).(*Animal1)) {
		t.Error("两个对象不是同一个")
	}
}

// 字段是接口
type Animal1 struct {
	Dog IAnimal `injection:"true"`
}

// 字段是结构体
type Animal2 struct {
	Dog Dog `injection:"true"`
}

// 字段是结构体指针
type Animal3 struct {
	Dog *Dog `injection:"true"`
}

func TestFieldType(t *testing.T) {
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()

	// 创建容器
	sc := &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	sc.AddServiceOf(goioc.Scope, imy, my)
	sc.AddService(goioc.Scope, my)

	t1 := reflect.TypeOf((*Animal1)(nil)).Elem()
	t2 := reflect.TypeOf((*Animal2)(nil)).Elem()
	t3 := reflect.TypeOf((*Animal3)(nil)).Elem()

	sc.AddService(goioc.Scope, t1)
	sc.AddService(goioc.Scope, t2)
	sc.AddService(goioc.Scope, t3)

	// 构建服务 Provider
	p := sc.Build()

	v1, _ := p.GetService(t1)
	v2, _ := p.GetService(t2)
	v3, _ := p.GetService(t3)

	fmt.Println(*v1)
	fmt.Println(*v2)
	fmt.Println(*v3)
}
