package socket

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/dgf/protzi/api/web"
	"github.com/gorilla/websocket"
)

type Client interface {
	Interrupt()
	Components()
	Flows()
	Flow(string)
	Add(flow, name, component string)
	Connect(flow, outName, outPort, inName, inPort string)
	Send(flow, name, port string, value interface{})
	Receive(flow, name, port string)
}

type client struct {
	conn     *websocket.Conn
	dispatch func(response web.Message)
}

func (c *client) send(request web.Command) {
	if err := c.conn.WriteJSON(web.CommandMessage(request)); err != nil {
		c.dispatch(*web.ErrorMessage(err))
	}
}

func (c *client) Components() {
	c.send(web.ComponentsCommand())
}

func (c *client) Flows() {
	c.send(web.FlowsCommand())
}

func (c *client) Flow(name string) {
	c.send(web.FlowCommand(name))
}

func (c *client) Add(flow, name, component string) {
	c.send(web.AddCommand(flow, name, component))
}

func (c *client) Connect(flow, outName, outPort, inName, inPort string) {
	c.send(web.ConnectCommand(flow, outName, outPort, inName, inPort))
}

func (c *client) Receive(flow, name, port string) {
	c.send(web.OutputCommand(flow, name, port))
}

func (c *client) Send(flow, name, port string, value interface{}) {
	if data, err := json.Marshal(value); err != nil {
		c.dispatch(*web.ErrorMessage(err))
	} else {
		c.send(web.InputCommand(flow, name, port, json.RawMessage(data)))
	}
}

func (c *client) Interrupt() {
	defer c.conn.Close()
	message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Tschüssi")
	if err := c.conn.WriteMessage(websocket.CloseMessage, message); err != nil {
		log.Println("interrupt close failed:", err)
		return
	}
}

func Connect(addr string, dispatch func(response web.Message)) Client {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c := &client{conn: conn, dispatch: dispatch}
	go func() {
		for {
			var response web.Message
			if err := c.conn.ReadJSON(&response); err != nil {
				log.Println("read failed:", err)
				return
			}
			dispatch(response)
		}
	}()
	return c
}
