package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func byteCount2(n int) string {
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

func home(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	start := time.Now()

	// Without limiting
	n, err := io.Copy(io.Discard, r.Body)

	// With limiting
	//n, err := io.Copy(io.Discard, io.LimitReader(r.Body, 100_000))

	if err != nil {
		log.Println("io copy:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	readLen := byteCount2(int(n))
	duration := time.Since(start)
	log.Printf("%s in %v\n", readLen, duration)

	fmt.Fprintf(w, "%s read in %v", readLen, duration)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	srv := &http.Server{
		Addr:    ":3000",
		Handler: mux,
		//ReadTimeout: 1 * time.Second,
		//WriteTimeout: 10 * time.Second,
		//IdleTimeout:  1 * time.Minute,
	}

	fmt.Println("starting server on port 3000...")
	log.Fatal(srv.ListenAndServe())
}
