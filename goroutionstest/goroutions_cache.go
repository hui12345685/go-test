package goroutionstest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func NewMemo(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)

		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	return e.res.value, e.res.err
}

func httpGetBody(url string) (interface{}, error) {
	//fmt.Printf("http get url:%s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func CacheTest() {
	fmt.Println("\033[1;32;40m  \nstart CacheTest---------------------------------- \033[0m")
	incomingURLs := []string{
		"https://golang.org", "https://godoc.org",
		"https://play.golang.org", "http://gopl.io",
		"https://golang.org", "https://godoc.org",
		"https://play.golang.org", "http://gopl.io",
	}
	m := NewMemo(httpGetBody)
	for _, url := range incomingURLs {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////
// 这个例子的channel使用真的是够牛逼了，正常哪里会想到这么复杂的通路逻辑
// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

type MemoV2 struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func NewV2(f Func) *MemoV2 {
	memo := &MemoV2{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *MemoV2) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *MemoV2) Close() { close(memo.requests) }

func (memo *MemoV2) server(f Func) {
	cache := make(map[string]*entry)
	//memo.requests channel感觉很多时候当成了队列来用，但是不知道性能如何
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition.
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

func CacheTestV2() {
	fmt.Println("\033[1;32;40m  \nstart CacheTestV2---------------------------------- \033[0m")
	incomingURLs := []string{
		"https://golang.org", "https://godoc.org",
		"https://play.golang.org", "http://gopl.io",
		"https://golang.org", "https://godoc.org",
		"https://play.golang.org", "http://gopl.io",
	}
	m := NewV2(httpGetBody)
	for _, url := range incomingURLs {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			log.Print(err)
		}
		fmt.Printf("%s, %s, %d bytes\n",
			url, time.Since(start), len(value.([]byte)))
	}
}
