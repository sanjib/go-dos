package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func byteCount(n int) string {
	unit := 1000
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	exp := -1
	for n >= unit {
		n /= unit
		exp++
	}
	return fmt.Sprintf("%d%s", n, string("KMGT"[exp]))
}

type slowloris struct {
	msg string
}

// Read simulates a slow read
func (s *slowloris) Read(p []byte) (int, error) {
	i := 0
	for i < len(p) {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(string(s.msg[i]))

		p = append(p, s.msg[i])
		i++

		if i == len(s.msg) {
			fmt.Println()
			return i, io.EOF
		}
	}
	return i, nil
}

// Read simulates read forever
//func (s *slowloris) Read(p []byte) (int, error) {
//	i := 0
//	for i < len(s.msg) {
//		time.Sleep(10 * time.Millisecond)
//		fmt.Print(string(s.msg[i]))
//
//		p = append(p, s.msg[i])
//		i++
//	}
//	return i, nil
//}

func postSlow(msg string) {
	start := time.Now()

	r := &slowloris{msg: msg}
	bufLen := byteCount(len(r.msg))

	err := post(r)
	if err != nil {
		log.Println(err)
		return
	}

	duration := time.Since(start)
	log.Printf("sent: %s in %v\n", bufLen, duration)
}

func postLarge(n int) {
	start := time.Now()

	var buf bytes.Buffer
	buf.Write(make([]byte, n))
	bufLen := byteCount(buf.Len())

	err := post(&buf)
	if err != nil {
		log.Println(err)
		return
	}

	duration := time.Since(start)
	log.Printf("sent: %s in %v\n", bufLen, duration)
}

func post(r io.Reader) error {
	resp, err := http.Post("http://localhost:3000/", "text/plain", r)
	if err != nil {
		return fmt.Errorf("post: %s", err)
	}
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read all: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp: %s", bb)
	}
	log.Printf("resp: %s\n", bb)
	return nil
}

func main() {
	// Large file
	postLarge(1_000_000_000)

	// Slowloris
	postSlow("I've seen things you people wouldn't believe... Attack ships on fire off the shoulder of Orion... I watched C-beams glitter in the dark near the Tannhauser Gate. All those moments will be lost in time, like tears in rain...")
}
