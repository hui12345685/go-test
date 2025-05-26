package goroutionstest

import (
	"fmt"
	"log"
	"mytest/closuretest"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawlV2(url string) []string {
	fmt.Printf("search url v2:%s\n", url)
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := closuretest.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

// 爬虫并发版本,带channel 和 goroutions的
func TestPageV2() {
	fmt.Println("\033[1;32;40m  \nstart TestPageV2---------------------------------- \033[0m")
	worklist := make(chan []string)
	n := 0 // number of pending sends to worklist
	n++
	go func() { worklist <- []string{"https://golang.org"} }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		//channel 收到数据之后才进来
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawlV2(link)
				}(link)
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////////////

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			info, _ := entry.Info()
			fileSizes <- info.Size()
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

func DuTest() {
	fmt.Println("\033[1;32;40m  \nstart DuTest---------------------------------- \033[0m")
	//roots := []string{"/data/home/yonghuiyu/work/qidian"}
	roots := []string{"E:\\work\\tencent\\qidian\\go", "E:\\work\\tencent\\qidian\\qd_cc"}

	fmt.Printf("du %+v\n", roots)
	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.2f GB\n", nfiles, float64(nbytes)/1e9)
}

// ////////////////////////////////////////////////////////////////////////////////////////
func DuTestV2() {
	fmt.Println("\033[1;32;40m  \nstart DuTestV2---------------------------------- \033[0m")
	//roots := []string{"/data/home/yonghuiyu"}
	roots := []string{"E:\\work\\tencent\\qidian\\go", "E:\\work\\tencent\\qidian\\qd_cc"}
	fmt.Printf("du %+v\n", roots)
	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

// ////////////////////////////////////////////////////////////////////////////////////////

type TestFileInfo struct {
	fileName string
	fileSize int64
}

func DuTestV3() {
	fmt.Println("\033[1;32;40m  \nstart DuTestV3---------------------------------- \033[0m")
	roots := []string{"/"}

	fmt.Printf("du %+v\n", roots)

	// Traverse the file tree.
	fileInfos := make(chan TestFileInfo)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDirV3(root, &n, fileInfos)
	}
	go func() {
		n.Wait()
		close(fileInfos)
		fmt.Printf("close file info:%d\n", fileInfos)
	}()

	// Print the results periodically.
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64
loop:
	for {
		select {
		case fileInfo, ok := <-fileInfos:
			if !ok {
				//关闭channle之后，ok会返回false
				fmt.Printf("DuTestV3 finished.\n")
				break loop // fileSizes was closed
			}
			if fileInfo.fileSize == 0 {
				continue
			}
			nfiles++
			nbytes += fileInfo.fileSize
			//fmt.Printf("read file:%+v,total files:%d,total size:%d\n", fileInfo, nfiles, nbytes)
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func walkDirV3(dir string, n *sync.WaitGroup, fileInfos chan<- TestFileInfo) {
	defer n.Done()
	for _, entry := range direntsV3(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDirV3(subdir, n, fileInfos)
		} else {
			info, err := entry.Info()
			fileInfo := TestFileInfo{}
			if info == nil || err != nil {
				//fmt.Printf("read failed, dir:%+v, error %s\n", entry, err)
				fileInfo.fileSize = 0
				fileInfos <- fileInfo
			} else {
				fileInfo.fileName = filepath.Join(dir, entry.Name())
				if fileInfo.fileName == "/proc/kcore" {
					//虚拟机的跳过，影响统计
					fileInfo.fileSize = 0
				} else {
					fileInfo.fileSize = info.Size()
				}
				fileInfos <- fileInfo
			}
		}
	}
}

// sema is a counting semaphore for limiting concurrency in dirents.
// 原来的代码channel的大小是20，这里改为3,降低cpu使用率
var sema = make(chan struct{}, 3)

// dirents returns the entries of directory dir.
func direntsV3(dir string) []os.DirEntry {
	//这里是有缓冲区的channel，只有到达20个才会阻塞
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, " dir:%s,read failed: %+v\n", dir, err)
		return nil
	}
	//fmt.Printf("read dir %s\n", dir)
	return entries
}
