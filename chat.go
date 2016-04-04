package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"log"
	"sync"
	"time"
)

// Should look like path
const websocketRoom = "/chat"

var connections = make(map[string]socketio.Socket)

func NewChatServer() (*socketio.Server, error) {
	lastMessages := []string{}
	var lmMutex sync.Mutex

	sio, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sio.On("connection", func(so socketio.Socket) {
		username := "User-" + so.Id()
		log.Println("number of connections: ", sio.Count())
		log.Println("on connection", username)
		so.Join(websocketRoom)

		lmMutex.Lock()
		for i, _ := range lastMessages {
			so.Emit("message", lastMessages[i])
		}
		lmMutex.Unlock()

		so.On("joined_message", func(message string) {
			username = message
			connections[username] = so
			log.Println("joined_message", message)
			res := map[string]interface{}{
				"username": username,
				"dateTime": time.Now().UTC(),
				"type":     "joined_message",
			}
			jsonRes, _ := json.Marshal(res)
			so.Emit("message", string(jsonRes))
			so.BroadcastTo(websocketRoom, "message", string(jsonRes))
		})

		so.On("send_message", func(message string) {
			log.Println("send_message from", username)
			res := map[string]interface{}{
				"username": username,
				"message":  message,
				"dateTime": time.Now().UTC(),
				"type":     "message",
			}
			jsonRes, _ := json.Marshal(res)
			lmMutex.Lock()
			if len(lastMessages) == 100 {
				lastMessages = lastMessages[1:100]
			}
			lastMessages = append(lastMessages, string(jsonRes))
			lmMutex.Unlock()
			so.Emit("message", string(jsonRes))
			so.BroadcastTo(websocketRoom, "message", string(jsonRes))
		})

		so.On("direct_message", func(data string) {
			log.Println("send direct message to a user", data)
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(data), &dat); err != nil {
				log.Println(err)
			}
			res := map[string]interface{}{
				"username": username,
				"message":  dat["message"].(string),
				"dateTime": time.Now().UTC(),
				"type":     "message",
			}
			jsonRes, _ := json.Marshal(res)
			so.Emit("message", string(jsonRes))
			if connections[dat["to"].(string)] != nil {
				connections[dat["to"].(string)].Emit("message", string(jsonRes))
			} else {
				log.Println(dat["to"], "is not existed")
			}
		})

		so.On("disconnection", func() {
			delete(connections, username)
			log.Println("on disconnect", username)
		})
	})

	sio.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return sio, err
}
