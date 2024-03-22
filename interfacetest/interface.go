package interfacetest

import "fmt"

type Base1 interface {
	Get(a int, b int) (int, error)
}

type Base2 interface {
	GetWithOption(a int, b int) (int, error)
}

type Drived1 struct {
}

func (*Drived1) Get(a int, b int) (int, error) {
	fmt.Println("Drived1 Get.")
	return 0, nil
}

func (*Drived1) GetWithOption(a int, b int) (int, error) {
	fmt.Println("Drived1 GetWithOption.")
	return 0, nil
}

func InterfaceConvertTest() {
	fmt.Println("\033[1;32;40m  \nstart ConvertTest---------------------------------- \033[0m")
	var d Drived1
	var base Base1
	base = &d
	//有点类似c++的多继承，预期结果也是一样的
	if _, ok := base.(Base1); ok {
		fmt.Println("get base1 type is ok.")
	} else {
		fmt.Println("get base1 type is failed.")
	}

	if _, ok := base.(Base2); ok {
		fmt.Println("get base2 type is ok.")
	} else {
		fmt.Println("get base2 type is failed.")
	}
}
