package closuretest

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func increase() func(int) int {
	sum := 0
	return func(i int) int {
		sum += i
		return sum
	}
}

func func1() {
	a := 1
	//defer延迟调用的时候还是保存的初始值，所以还是1
	defer func(r int) {
		fmt.Println(r)
	}(a)
	a = a + 100
	fmt.Println(a)
}

func func2() {
	num := 0
	for i := 0; i < 5; i++ {
		go func() {
			//这里多协程就和多线程一样，各个协程里面看到num的值完全是不确定的(0-4之间)
			fmt.Println(num)
			num++
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(num)
}

func Test1() {
	fmt.Println("\033[1;32;40m  \nstart closure.Test1---------------------------------- \033[0m")
	incr := increase()
	fmt.Println(incr(1))
	fmt.Println(incr(2))

	func1()
	func2()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// 装饰器模式
// 为函数类型设置别名提高代码可读性
type MultiPlyFunc func(int, int) int

// 乘法运算函数1（算术运算）
func multiply1(a, b int) int {
	return a * b
}

// 乘法运算函数2（位运算）
func multiply2(a, b int) int {
	return a << b
}

// 通过高阶函数在不侵入原有函数实现的前提下计算乘法函数执行时间
func execTime(f MultiPlyFunc) MultiPlyFunc {
	return func(a, b int) int {
		start := time.Now()      // 起始时间
		c := f(a, b)             // 执行乘法运算函数
		end := time.Since(start) // 函数执行完毕耗时
		fmt.Printf("--- 执行耗时: %v ---\n", end)
		return c // 返回计算结果
	}
}

func Test2() {
	fmt.Println("\033[1;32;40m  \nstart closure.Test2---------------------------------- \033[0m")
	a := 2
	b := 8
	fmt.Println("算术运算：")
	decorator1 := execTime(multiply1)
	c := decorator1(a, b)
	fmt.Printf("%d x %d = %d\n", a, b, c)
	fmt.Println("位运算：")
	decorator2 := execTime(multiply2)
	a = 1
	b = 4
	c = decorator2(a, b)
	fmt.Printf("%d << %d = %d\n", a, b, c)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// 这里是一个类似爬虫的代码，从一个url出发，找到所有关联的url，然后依次往后搜索url，类似于广度搜索所有关联的url
// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	fmt.Printf("get links %+v \nfrom url:%s\n", links, url)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	//这里是类似一个树形的扩散
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
			fmt.Printf("breadthFirst url:%s\n", item)
		}
	}
}

func crawl(url string) []string {
	fmt.Printf("search url:%s\n", url)
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func TestPage() {
	fmt.Println("\033[1;32;40m  \nstart closure.TestPage---------------------------------- \033[0m")
	// pages := []string{"https://golang.org"}
	pages := []string{"https://www.baidu.com"}
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, pages)
}
