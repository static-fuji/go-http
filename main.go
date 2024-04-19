package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
)

var (
	data  map[string][]byte
	mutex sync.RWMutex
)

func main() {
	http.HandleFunc("/objects/", handler)
	http.ListenAndServe(":8000", nil)

	fmt.Println("サーバーを起動します... ポート 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("サーバーの起動に失敗しました:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/objects/"):]

	validKey := regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)
	if !validKey.MatchString(key) {
		http.Error(w, "404, Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPut:
		putData(w, r, key)
	case http.MethodGet:
		getData(w, r, key)
	default:
		http.Error(w, "405, Method not allowed", http.StatusMethodNotAllowed)
	}

	fmt.Fprintf(w, "アクセスキー:%s", key)
}

func putData(w http.ResponseWriter, r *http.Request, key string) {
	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	data[key] = body
}

func getData(w http.ResponseWriter, r *http.Request, key string) {
	mutex.RLock()
	defer mutex.RUnlock()

	// データを取得
	data, ok := data[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// レスポンスにデータを書き込み
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
