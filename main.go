package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type Mensaje struct {
	Nombre string
	Msg    string
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	var broadcast = socketio.NewBroadcast()
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		broadcast.Join("chat_room", s)
		return nil
	})
	server.OnEvent("/", "reply", func(s socketio.Conn, msg string) {
		fmt.Println("reply:", msg)
		var mensaje = []byte(msg)
		//s.Emit("reply", msg)
		var obj Mensaje
		err := json.Unmarshal(mensaje, &obj)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Nombre:" + obj.Nombre + " mensaje: " + obj.Msg)
		broadcast.Send("chat_room", "reply", msg, func(so socketio.Broadcast, data string) {
			log.Println("Client ACK with data: ", data)
		})

	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
