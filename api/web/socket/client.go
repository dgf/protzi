package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client interface {
	Interrupt()
	Send(Command) error
}

type client struct {
	conn *websocket.Conn
}

func (c *client) Interrupt() {
	defer c.conn.Close()
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Tschüssi")
	if err := c.conn.WriteMessage(websocket.CloseMessage, message); err != nil {
		log.Println("interrupt close failed:", err)
		return
	}
}

func (c *client) Send(request Command) error {
	if payload, err := json.Marshal(request); err != nil {
		return err
	} else {
		message := []byte(fmt.Sprintf(`{"type":"command","payload":%s}`, payload))
		return c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func Connect(addr string, dispatch func(response Message)) Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c := &client{conn: conn}
	go func() {
		for {
			var response Message
			if err := c.conn.ReadJSON(&response); err != nil {
				log.Println("read failed:", err)
				return
			}
			dispatch(response)
		}
	}()
	return c
}
