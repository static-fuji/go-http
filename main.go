package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

// mutexを含むデータ
type dataBody struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

// データのインスタンス
func New() *dataBody {
	return &dataBody{
		data: make(map[string][]byte),
	}
}

func main() {
	m := New()
	http.HandleFunc("/objects/", m.handler)
	http.ListenAndServe(":8000", nil)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("サーバーの起動に失敗しました:", err)
	}
}

// リクエストの管理
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

// リクエストヘッダがPUTのとき
func (m *dataBody) putData(w http.ResponseWriter, r *http.Request, key string) {
	//リクエストボディの読み込み
	body, err := io.ReadAll(r.Body)

	//エラー処理
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	//読み込んだデータの保持
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = body
}

// リクエストヘッダがGETのとき
func (m *dataBody) getData(w http.ResponseWriter, r *http.Request, key string) {

	//レスポンスボディの作成
	m.mutex.RLock()
	data, ok := m.data[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	//レスポンスの書き込み
	m.mutex.RUnlock()
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
