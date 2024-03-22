package cgotest

/*
#include "testc.h"
#include "cwrap.h"
*/
import "C"

import "fmt"

// go调用c确实是很方便的
// 如果是调用编译好的库文件也不麻烦，在上面的注释代码中加上#cgo CFLAGS:/#cgo LDFLAGS:,网上有一大堆资料
func GoSumTest(a, b int) int {
	fmt.Println("\033[1;32;40m  \nstart c test GoSumTest---------------------------------- \033[0m")
	s := C.Sum(C.int(a), C.int(b))
	fmt.Println(s)
	return a + b
}

// go调用c++确实是很方便的
// 如果是调用编译好的库文件也不麻烦，在上面的注释代码中加上#cgo CFLAGS:/#cgo LDFLAGS:,网上有一大堆资料
func GoCppTest() {
	fmt.Println("\033[1;32;40m  \nstart cpp test GoCppTest---------------------------------- \033[0m")
	C.call()
}
