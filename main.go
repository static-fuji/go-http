package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/objects/"):]

	validKey := regexp.MustCompile(`^[a-zA-Z0-9]{1,10}$`)
	if !validKey.MatchString(key) {
		http.Error(w, "404", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "アクセスキー:%s", key)
}

func main() {
	http.HandleFunc("/objects/", handler)
	http.ListenAndServe(":8000", nil)

	fmt.Println("サーバーを起動します... ポート 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("サーバーの起動に失敗しました:", err)
	}
}
