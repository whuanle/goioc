# goioc

goioc 是一个基于 GO 语言编写的依赖注入框架，代码量不多，很简洁。

```bash
go get -u https://github.com/whuanle/goioc
```



### 接口注入

有以下接口定义：

```golang
type IAnimal interface {
	Println(s string)
}
```

接口的实现：

```golang
type Dog struct {
}
func (my Dog) Println(s string) {
	fmt.Println(s)
}
```



当使用依赖注入框架时，我们可以将接口和实现分开，放到两个模块中。



首先创建 IServiceCollection 容器

```go
	// 创建容器
	var collection IServiceCollection = &ServiceCollection{}
```



获取接口和结构体的 reflect.Type：

```go

// 写法 1
    // 接口的 reflect.Type
	var animal IAnimal
    imy := reflect.TypeOf(&animal).Elem()
	my := reflect.TypeOf(Dog{})

// 写法 2
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()
```

> 以上两种写法都可以使用，目的在于获取到接口和结构体的 reflect.Type。不过第一种方式会实例化结构体，消耗了一次内存，并且要获取接口的 reflect.Type，是不能直接有用 `reflect.TypeOf(animal)` 的，需要使用 `reflect.TypeOf(&animal).Elem()` 。



然后注入服务，其生命周期为 Scoped：

```go
	// 注入服务，生命周期为 scoped
	collection.AddScopedForm(imy, my)
```

> 当你需要 IAnimal 接口时，会自动注入 Dog 结构体给 IAnimal。



构建依赖注入服务提供器：

```go
	// 构建服务 Provider
	serviceProvider := collection.Build()
```



构建完成后，即可通过 Provider 对象获取需要的实例：

```go
	// 获取对象
	// *interface{}
	obj, err := serviceProvider.GetService(imy)
	if err != nil {
		panic(err)
	}
	
	// 转换为接口
	a := (*obj).(IAnimal)
	// 	a := (*obj).(*Dog)
```

因为使用了依赖注入，我们使用时，只需要使用接口即可，不需要知道具体的实现。



完整的代码示例：

```go
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
	// *interface{} = &Dog{}
	obj, err := serviceProvider.GetService(imy)

	if err != nil {
		panic(err)
	}

	fmt.Println("obj 类型是", reflect.ValueOf(obj).Type())

	// *interface{} = &Dog{}，因此需要处理指针
	animal := (*obj).(IAnimal)
	// 	a := (*obj).(*Dog)
	animal.Println("测试")
```






### 结构体字段依赖注入

结构体中的字段，可以自动注入和转换实例。

如结构体 Animal 的定义，其使用了其它结构体，goioc 可以自动注入 Animal 对应字段。要被注入的字段必须是接口或者结构体！

```go
type IAnimal interface {
	Println(s string)
}

type Dog struct {
}
func (my Dog) Println(s string) {
	fmt.Println(s)
}

// 结构体中包含了其它对象
type Animal struct {
	Dog IAnimal `injection:"true"`
}
```

> 要对需要自动注入的字段设置 tag 中包含`injection:"true"` 才会生效。

依赖注入的实例代码如下：

```go
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
```



goioc 可以自动给你的结构体字段进行自动依赖注入。

> 注意，如果 obj 要转换为接口，则是使用：
>
> ```
> 	animal := (*obj).(IAnimal)
> ```
>
> 如果 obj 要转换为结构体，则是：
>
> ```
> 	animal := (*obj).(*Animal)
> ```





### 生命周期

对象的生命周期有三个：

```go
const (
	Transient ServiceLifetime= iota
	Scope
	Singleton
)
```

如果是单例模式，则在同一个容器中，无论多少次构建 Provider，以及使用 GetService 获取对象，每次获取到的都是同一个对象。

如果是 Scope 模式，ServiceCollection 每次 Build 时，同一个 serviceProvider 获取到的对象是同一个。

如果是 Transient 模式，每次获取到的对象都是新的。

示例：

```go
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()
	var collection IServiceCollection = &ServiceCollection{}
	collection.AddScopedForm(imy, my)
	serviceProvider := collection.Build()

	// 获取对象
	// *interface{} = &Dog{}，因此需要处理指针
	obj1, _ := serviceProvider.GetService(imy)
	obj2, _ := serviceProvider.GetService(imy)

	fmt.Printf("obj1 = %p,obj2 = %p", (*obj1).(*Dog), (*obj2).(*Dog))
	if fmt.Sprintf("%p",(*obj1).(*Dog)) != fmt.Sprintf("%p",(*obj2).(*Dog)){
		t.Error("两个对象不是同一个")
	}
```





### 自定义实例化过程

默认情况下，注入的结构体是直接实例化的，如：

```go
dog := &Dog{
}
```

如果结构体有字段需要注入，则：

```go
animal := &Animal{
		Dog: &Dog{
		},
	}
```



如果你想自定义实例化过程，则可以注入一个匿名函数：

```go
	// 获取 reflect.Type
	imy := reflect.TypeOf((*IAnimal)(nil)).Elem()
	my := reflect.TypeOf((*Dog)(nil)).Elem()

	// 创建容器
	var collection IServiceCollection = &ServiceCollection{}

	// 注入服务，生命周期为 scoped
	collection.AddScopedHandlerForm(imy, my,func()interface{}{
		// 你自己的实例化代码
		return &Dog{
		}
	})
```



