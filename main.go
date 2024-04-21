package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type dataBody struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

func New() *dataBody {
	return &dataBody{
		data: make(map[string][]byte),
	}
}

func main() {
	m := New()
	http.HandleFunc("/objects/", m.handler)
	http.ListenAndServe(":8000", nil)

	fmt.Println("サーバーを起動します... ポート 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("サーバーの起動に失敗しました:", err)
	}
}

func (m *dataBody) handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/objects/"):]

	validKey := regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)
	if !validKey.MatchString(key) {
		http.Error(w, "404, Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPut:
		m.putData(w, r, key)
	case http.MethodGet:
		m.getData(w, r, key)
	default:
		http.Error(w, "405, Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (m *dataBody) putData(w http.ResponseWriter, r *http.Request, key string) {
	//body := make([]byte, r.ContentLength)
	fmt.Println(m.data[key])
	body, err := io.ReadAll(r.Body)

	//mapの初期化

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	//data := make(map[string][]byte)

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = body
	fmt.Println(m.data[key])
}

func (m *dataBody) getData(w http.ResponseWriter, r *http.Request, key string) {
	fmt.Println(m.data[key])

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// データを取得
	data, ok := m.data[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// レスポンスにデータを書き込み
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
