package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":3030")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() {
		// NOTE: ignoring errors.
		_, _ = io.Copy(os.Stdout, conn)

		log.Println("done")
		// signal the main goroutine.
		done <- struct{}{}
	}()

	mustCopy(conn, os.Stdin)
	_ = conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
