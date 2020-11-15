package handler

import (
	"fmt"
	"net/http"
)

// HttpInterceptor http请求拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			token := r.Header.Get("token")
			fmt.Println(token)

			h(w, r)
		})
}
