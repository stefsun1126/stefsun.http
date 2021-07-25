package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	/*
		開啟61個goroutine，須從1打印到60，且不能出現重複，大於60打印Error
	*/

	for i := 0; i < 61; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get("http://127.0.0.1:8080")
			if err != nil {
				log.Fatal("get failed, err:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("read from resp.Body failed, err:", err)
				return
			}
			fmt.Println("body", string(body))
		}()
	}
	wg.Wait()
}
