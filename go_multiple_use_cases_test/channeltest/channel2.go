package channeltest

import (
	"fmt"
)

// /////////////////////////////////////////////////////////////////////
func TestCounter() {
	fmt.Println("\033[1;32;40m  \nstart TestCounter---------------------------------- \033[0m")
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; x < 10; x++ {
			fmt.Printf("write before:%d\n", x)
			naturals <- x
			fmt.Printf("write after:%d\n", x)
		}
		close(naturals)
	}()

	// Squarer
	go func() {
		//使用range的方式遍历，貌似不会完全阻塞channel
		for x := range naturals {
			fmt.Printf("read after:%d\n", x)
			squares <- x * x
		}
		close(squares)
	}()

	// Printer (in main goroutine)
	for x := range squares {
		fmt.Println(x)
	}
}

// //////////////////////////////////////////////////////////////////////////
// 和上面的case是一样的，只是吧goroution函数抽出来了
func counter(out chan<- int) {
	for x := 0; x < 10; x++ {
		fmt.Printf("v2 write:%d\n", x)
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		fmt.Printf("v2 read:%d\n", v)
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

func TestCounterV2() {
	fmt.Println("\033[1;32;40m  \nstart TestCounterV2---------------------------------- \033[0m")
	naturals := make(chan int)
	squares := make(chan int)
	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

// //////////////////////////////////////////////////////////////////////////
// select 下面只会有一个分支会执行
// 如果多个case同时就绪时，select会随机地选择一个执行，这样来保证每一个channel都有平等的被select的机会
func SelectTest() {
	fmt.Println("\033[1;32;40m  \nstart SelectTest---------------------------------- \033[0m")
	ch := make(chan int, 1) //缓冲区为1，会先写channel成功
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Println(x) // "0" "2" "4" "6" "8"
		case ch <- i:
		}
	}
}
