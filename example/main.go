package main

import (
	"fmt"
	ioc "github.com/whuanle/goioc"
)

type IAnimal interface {
	Println(s string)
}

type Dog struct {
}

func (my Dog) Println(s string) {
	fmt.Println(s)
}
func main() {
	// 创建容器
	var collection ioc.IServiceCollection = &ioc.ServiceCollection{}
	ioc.AddScopedForm[IAnimal, Dog](collection)
	provider := collection.Build()
	v := ioc.GetI[IAnimal](provider)
	v.Println("a")
}
