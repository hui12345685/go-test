package servertest

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mu    sync.Mutex
	count int
)

func TestHttpSvrV1() {
	fmt.Println("\033[1;32;40m  \nstart TestHttpSvrV1---------------------------------- \033[0m")

	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	addr := "localhost:8000"
	//addr := "127.0.0.1:8000"
	log.Fatal(http.ListenAndServe(addr, nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle request: %+v\n", *r)
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("counter request: %+v\n", *r)
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

//////////////////////////////////////////////////////////////////////////////////////////////////

// Clock1 is a TCP server that periodically writes the time.

func TestTcpSvrV1() {
	fmt.Println("\033[1;32;40m  \nstart TestTcpSvrV1---------------------------------- \033[0m")

	listener, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	n := 0
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
		n++
		if n >= 10 {
			break
		}
	}
}
