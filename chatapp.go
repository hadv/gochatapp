package main

import (
	"net/http"
	"runtime"
)

func main() {
	// Sets the number of maxium goroutines to the 2*numberCPU + 1
	runtime.GOMAXPROCS((runtime.NumCPU() * 2) + 1)

	// Sets up the handlers and listen on port 8080
	chat, _ := NewChatServer()
	http.Handle("/socket.io/", chat)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", http.FileServer(http.Dir("./templates/")))
	http.ListenAndServe(":8888", nil)
}
