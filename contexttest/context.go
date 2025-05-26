package contexttest

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Println("hello")
			time.Sleep(70 * time.Millisecond)
		case <-ctx.Done():
			fmt.Printf("will exited, reason:%+v\n", ctx.Err())
			return ctx.Err()
		}
	}
}

// 并发体超时或 ContextTest 主动停止工作者 Goroutine 时，每个工作者都可以安全退出
func ContextTest() {
	fmt.Println("\033[1;32;40m  \n start ContextTest--------------------------------------- \033[0m")
	//第二个参数时间表示goroution经过1秒超时后，ctx.Done()这个管道会被写入数据，ctx.Err()为context deadline exceeded
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(ctx, &wg)
	}
	//如果这个sleep设置的时间大于上面WithTimeout的时间，会先超时结束
	//time.Sleep(100 * time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	fmt.Println("will cancel goroution")
	//调用cancel之后ctx.Done()这个管道会被写入数据，然后上面的worker函数对应的分支就会成立,ctx.Err()为context canceled
	cancel()

	wg.Wait()
}

/////////////////////////////////////////////////////////////////////////////////////

// 返回生成自然数序列的管道: 2, 3, 4, ...
func GenerateNatural(ctx context.Context, wg *sync.WaitGroup) chan int {
	ch := make(chan int)
	go func() {
		defer fmt.Printf("exit GenerateNatural goroutine,ch %+v\n", ch)
		fmt.Printf("enter GenerateNatural goroutine,ch %+v\n", ch)
		defer wg.Done()
		defer close(ch)
		for i := 2; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				fmt.Printf("GenerateNatural write channel data:%d,ch %+v\n", i, ch)
			}
			// fmt.Printf("GenerateNatural: idx %v\n", i)
		}
	}()
	return ch
}

// 管道过滤器: 删除能被素数整除的数
func PrimeFilter(ctx context.Context, in <-chan int, prime int, wg *sync.WaitGroup) chan int {
	out := make(chan int)
	go func() {
		defer wg.Done()
		defer close(out)
		defer fmt.Printf("exit PrimeFilter goroutine,prime:%d in %+v out %+v\n", prime, in, out)
		fmt.Printf("enter PrimeFilter goroutine,in %+v,prime:%d out %+v\n", in, prime, out)
		for i := range in {
			fmt.Printf("PrimeFilter read channel data:%d,prime:%d in %v out %v\n", i, prime, in, out)
			if i%prime != 0 {
				select {
				case <-ctx.Done():
					return
				case out <- i:
				}
				//fmt.Printf("not zero:idx %v prime:%v -- ", i, prime)
			}
			// fmt.Printf("Filter: idx %v\n", i)
		}

	}()
	return out
}

func ContextTest2() {
	fmt.Println("\033[1;32;40m  \n start ContextTest2--------------------------------------- \033[0m")
	wg := sync.WaitGroup{}
	// 通过 Context 控制后台 Goroutine 状态
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	ch := GenerateNatural(ctx, &wg) // 自然数序列: 2, 3, 4, ...
	fmt.Printf("get GenerateNatural ch addr:%+v\n", ch)
	for i := 0; i < 20; i++ {
		prime := <-ch // 新出现的素数
		fmt.Printf("main read channel data,%v: %v , ch addr:%+v\n", i+1, prime, ch)
		wg.Add(1)
		// PrimeFilter返回的ch是下一个协程的输入
		ch = PrimeFilter(ctx, ch, prime, &wg) // 基于新素数构造的过滤器
	}

	cancel()
	wg.Wait()
	fmt.Println("\033[1;32;40m  \n end ContextTest2--------------------------------------- \033[0m")
}
