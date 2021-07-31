package httppool

import "net/http"

func main() {
	var s server
	http.ListenAndServe("localhost:9999", &s)
}
