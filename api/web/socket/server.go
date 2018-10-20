package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dgf/protzi/api"
	"github.com/dgf/protzi/component/core"
	"github.com/dgf/protzi/component/text"
	"github.com/gorilla/websocket"
)

type Server interface {
	Start() error
	Endpoint() string
}

var upgrader = websocket.Upgrader{}

const (
	messageWait = 7 * time.Second
	writeWait   = 1 * time.Second
	pingPeriod  = messageWait - writeWait
)

func marshal(t string, v interface{}) *Message {
	if data, err := json.Marshal(v); err != nil {
		return Error(fmt.Errorf("%s marshal failed: %s", t, err.Error()))
	} else {
		return &Message{Type: t, Payload: json.RawMessage(data)}
	}
}

func flows() *Message {
	f := map[string]string{}
	for _, name := range api.Flows() {
		f[name] = "desc " + " " + name
	}
	return marshal("flows", f)
}

func components() *Message {
	c := map[string]string{}
	for _, name := range api.Components() {
		c[name] = "desc " + " " + name
	}
	return marshal("components", c)
}

func dispatch(request *Message) *Message {
	switch request.Type {
	case "command":
		var command Command
		if err := json.Unmarshal(request.Payload, &command); err != nil {
			return Error(fmt.Errorf("invalid command payload: %q", string(request.Payload)))
		}
		switch command.Call {
		case "flows":
			return flows()
		case "components":
			return components()
		default:
			return Error(fmt.Errorf("invalid command call: %q", command.Call))
		}
	default:
		return Error(fmt.Errorf("unknown dispatch: %q", string(request.Payload)))
	}
}

func handle(w http.ResponseWriter, r *http.Request) {

	// upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}
	defer conn.Close()

	// limit read
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(messageWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(messageWait))
		return nil
	})

	// read requests
	requests := make(chan Message)
	defer close(requests)
	out := make(chan struct{})
	defer close(out)
	go func() {
		for {
			var request Message
			if err := conn.ReadJSON(&request); err != nil {
				log.Println("read failed:", err)
				out <- struct{}{}
				return
			}
			requests <- request
		}
	}()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	for {
		select {
		case <-out:
			return
		case request := <-requests:
			response := dispatch(&request)
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(&response); err != nil {
				log.Println("write failed:", err)
				return
			}
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ping failed:", err)
				return
			}
		}
	}
}

type server struct {
	URL      string
	listener net.Listener
}

func Serve(addr string) Server {

	// register available components
	api.Register("Echo", &core.Echo{})
	api.Register("Print", &core.Print{})
	api.Register("Tick", &core.Tick{})
	api.Register("Time", &core.Time{})
	api.Register("Read", &text.FileRead{})
	api.Register("Render", &text.Render{})
	api.Register("WordCount", &text.WordCount{})

	// observer TCP
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	return &server{
		URL:      listener.Addr().(*net.TCPAddr).String(),
		listener: listener,
	}
}

func (s *server) Endpoint() string {
	return "ws://" + s.URL
}

func (s *server) Start() error {
	http.HandleFunc("/", handle)
	return http.Serve(s.listener, nil)
}
