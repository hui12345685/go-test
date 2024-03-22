package channeltest

import (
	"fmt"
	"sync"
	"time"
)

// ////////////////////////////////////////////////////////
func ChannelTest() {
	fmt.Println("\033[1;32;40m  \nstart ChannelTest---------------------------------------- \033[0m")
	c := make(chan int) //声明一个int类型的无缓冲通道
	go func() {
		fmt.Println("ready to send in g1")
		c <- 1
		fmt.Println("send 1 to chan")
		fmt.Println("goroutine start sleep 1 second")
		time.Sleep(time.Second)
		fmt.Println("goroutine end sleep")
		c <- 2
		fmt.Println("send 2 to chan")
	}()

	fmt.Println("main thread start sleep 1 second")
	time.Sleep(time.Second)
	fmt.Println("main thread end sleep")
	i := <-c
	fmt.Printf("receive %d\n", i)
	i = <-c
	fmt.Printf("receive %d\n", i)
	time.Sleep(time.Second)
}

// ////////////////////////////////////////////////////////
var (
	dog  = make(chan struct{})
	cat  = make(chan struct{})
	fish = make(chan struct{})
)

func Dog(n *sync.WaitGroup) {
	//n.Add(1)
	//defer n.Done()
	<-fish
	fmt.Println("dog")
	dog <- struct{}{}
}

func Cat(n *sync.WaitGroup) {
	//n.Add(1)
	//defer n.Done()
	<-dog
	fmt.Println("cat")
	cat <- struct{}{}
}

func Fish(n *sync.WaitGroup) {
	//n.Add(1)
	//defer n.Done()
	<-cat
	fmt.Println("fish")
	fish <- struct{}{}
	//fmt.Println("send to fish")
}

// 依次打印dog，cat，fish
func PrintCatFishDog() {
	fmt.Println("\033[1;32;40m  \nstart PrintCatFishDog-------------------------------- \033[0m")
	var n sync.WaitGroup
	for i := 0; i < 5; i++ {
		go Dog(&n)
		go Cat(&n)
		go Fish(&n)
	}
	fish <- struct{}{}
	fmt.Println("Wait for finish.")
	//sleep之后主goroution结束了，fish因为进程的退出而退出
	time.Sleep(time.Second * 1)
	//wait goroution stop
	//n.Wait() //用这种方式会卡死，因为发给fish没有其他goroution接收
}

// ////////////////////////////////////////////
// 注意这里缓冲区大小为1，所以可以在创建goroution之前写管道
var word = make(chan struct{}, 1)
var num = make(chan struct{}, 1)

func printNums() {
	for i := 0; i < 5; i++ {
		<-word
		fmt.Println(1)
		num <- struct{}{}
	}
}
func printWords() {
	for i := 0; i < 5; i++ {
		<-num
		fmt.Println("a")
		word <- struct{}{}
	}
}

// 依次打印a，1
func PrintWordAndNums() {
	fmt.Println("\033[1;32;40m  \nstart PrintWordAndNums---------------------------------- \033[0m")
	num <- struct{}{}
	go printNums()
	go printWords()
	time.Sleep(time.Second * 1)
}
