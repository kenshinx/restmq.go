package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"testing"
)

const (
	addr  = "http://127.0.0.1:8000"
	queue = "kenshin"
)

func BenchmarkPut(b *testing.B) {
	api := fmt.Sprintf("%s/q/%s", addr, queue)

	for i := 0; i < b.N; i++ {
		payload := url.Values{"value": []string{strconv.Itoa(i)}}
		http.PostForm(api, payload)
	}

}

func BenchmarkGet(b *testing.B) {
	api := fmt.Sprintf("%s/q/%s", addr, queue)

	for i := 0; i < b.N; i++ {
		http.Get(api)
	}
}

func isAccessed() bool {
	url := addr + "/q"
	_, err := http.Get(url)
	return err == nil
}

func init() {

	runtime.GOMAXPROCS(4)

	if !isAccessed() {
		fmt.Printf("Please make sure %s is the rest server address\n", addr)
		os.Exit(1)
	}
}
