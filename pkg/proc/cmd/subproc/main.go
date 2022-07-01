package main

import (
	"fmt"
	"time"

	"github.com/cybriq/glint/pkg/proc"
	"github.com/cybriq/qu"
)

func main() {
	quit := qu.T()
	p := proc.Consume(
		quit, func(b []byte) (e error) {
			fmt.Println("from child:", string(b))
			return
		}, "go", "run", "serve/main.go",
	)
	for {
		_, e := p.StdConn.Write([]byte("ping"))
		if e != nil {
			fmt.Println("err:", e)
		}
		time.Sleep(time.Second)
	}
}
