package main

import "net/http"

// 向客户端写数据
func wsHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World！"))
}

func main(){
	http.HandleFunc("/ws",wsHandler)

	_ = http.ListenAndServe("0.0.0.0:7777", nil)
}
