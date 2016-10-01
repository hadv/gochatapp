package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/zhouhui8915/go-socket.io-client"
	"net/http"
)

var _ = Describe("Chat", func() {
	var server *socketio.Server

	BeforeEach(func() {
		var err error
		if server, err = NewChatServer(); err != nil {
			fmt.Println("Error", err.Error())
			return
		}
		http.Handle("/socket.io/", server)

		go http.ListenAndServe(":8888", nil)
	})

	AfterEach(func() {
		fmt.Println("Test finished!!!")
	})

	Describe("join chat", func() {
		Context("join chat sucessfully", func() {
			var socket *socketio_client.Client
			var err error

			JustBeforeEach(func() {
				opts := &socketio_client.Options{
					Transport: "websocket",
					Query:     make(map[string]string),
				}
				socket, err = socketio_client.NewClient("http://localhost:8888", opts)
				if err != nil {
					Fail("cannot create socket client")
				}
			})

			It("join chat sucessfully", func(done Done) {
				socket.Emit("joined_message", "hadv")
				c := make(chan string)
				socket.On("message", func(msg string) {
					c <- msg
				})
				var dat map[string]interface{}

				if err := json.Unmarshal([]byte(<-c), &dat); err != nil {
					fmt.Println("invalid input data", err.Error())
					return
				}
				Expect(dat["type"]).To(Equal("joined_message"))
				Expect(dat["username"]).To(Equal("hadv"))
				close(done)
			}, 0.2)
		})
	})
})
