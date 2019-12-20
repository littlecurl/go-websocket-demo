package main

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		go func() {
			// 延迟调用（一种优雅的回收资源方式）
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					// handle error
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					// handle error
				}
			}
		}()
	}))
}
