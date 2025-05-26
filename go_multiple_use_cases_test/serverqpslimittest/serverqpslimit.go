package serverqpslimittest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// 令牌桶限流
type LimiteTest struct {
	limite      int
	tokenBucket chan struct{}
}

// 这里初始化，并启动一个定时的ticker，ticker.C 每10ms会写入一次数据
func (lt *LimiteTest) InitQpsLimiteTest(capacity int) {
	//把1秒分配1000份,所以限制的qps不能大于1000，如大于1000，则每次可以写多次channel
	//并且这里很不精确
	ms := 1000 / capacity
	var fillInterval = time.Millisecond * time.Duration(ms)
	lt.limite = capacity
	lt.tokenBucket = make(chan struct{}, lt.limite)

	fillToken := func() {
		ticker := time.NewTicker(fillInterval)
		for { //这里需要一个for，否则一次就退出去了
			select {
			case <-ticker.C:
				fmt.Println("will add some,current token cnt:", len(lt.tokenBucket), time.Now())
				select {
				case lt.tokenBucket <- struct{}{}:
				default:
				}
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	go fillToken()
	//time.Sleep(time.Hour)
	//select {}
}

// 每来一个请求从channel中取一条数据
func (lt *LimiteTest) TakeAvailable(block bool) bool {
	var takenResult bool
	if block {
		select {
		case <-lt.tokenBucket:
			takenResult = true
		}
	} else {
		select {
		case <-lt.tokenBucket:
			takenResult = true
		default:
			takenResult = false
		}
	}

	return takenResult
}

// 没有数据了，说明令牌用完了，需要被限制频率
func (lt *LimiteTest) IsLimited() bool {
	return len(lt.tokenBucket) == 0
}

// 全局变量，给单例返回的
var (
	limit *LimiteTest
	once  sync.Once
)

// 单例函数，为了使用Dao而暴露
func GetInst() *LimiteTest {
	once.Do(func() {
		limit = &LimiteTest{}
	})
	return limit
}

func sayhello(wr http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Millisecond) //压测工具太快了，限速一下
	limite := GetInst().IsLimited()
	if limite {
		fmt.Println("is limited.time:", time.Now())
		wr.WriteHeader(200)
		io.WriteString(wr, "is frequence limited")
		return
	}

	wr.WriteHeader(200)
	io.WriteString(wr, "hello world")
	GetInst().TakeAvailable(true) //如果超频了会阻塞
}

func SvrQpsLimiteTest() {
	GetInst().InitQpsLimiteTest(100)

	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
