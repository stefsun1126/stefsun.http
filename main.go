package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	configuration "stefsun.http/configuration"
	redis "stefsun.http/redis"
)

var (
	// 限制連接次數
	limitTimes int64 = 60
	// 限制時間
	limitDuration time.Duration = 60 * time.Second

	// 用來避免資源競爭
	lockChan = make(chan struct{}, 1)
)

// 將值取出
func getLock() {
	<-lockChan
}

// 塞值
func releaseLock() {
	lockChan <- struct{}{}
}

func main() {
	// 獲得config instance
	config := configuration.New("app.config", os.Getenv("MODE"))

	// 獲得redis instance
	rc, err := redis.NewRedis(config.GetString("service.redis.uri"))
	if err != nil {
		log.Fatal("NewRedis: ", err)
	}

	// 路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 沒有要用favicon 所以濾掉瀏覽器自己發起的請求
		if r.URL.RequestURI() == "/favicon.ico" {
			return
		}
		// 塞值到channle 如已經滿了 其他goroutine會阻塞
		releaseLock()
		// 根據redis內部的資料來判斷是否還可繼續訪問
		isAllow, currTimes := rc.AllowRequest(strings.Split(r.RemoteAddr, ":")[0], limitTimes, limitDuration)
		// 將channle值取出
		getLock()
		if isAllow {
			w.Write([]byte(strconv.Itoa(int(currTimes))))
		} else {
			w.Write([]byte("Error"))
		}
	})

	// 建立server監聽
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
