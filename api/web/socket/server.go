package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dgf/protzi/api"
	"github.com/dgf/protzi/api/web"
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

func marshal(t string, v interface{}) *web.Message {
	if data, err := json.Marshal(v); err != nil {
		return web.ErrorMessage(fmt.Errorf("marshal %q failed: %v", t, err))
	} else {
		return &web.Message{Type: t, Payload: json.RawMessage(data)}
	}
}
func forward(output api.Output, data web.OutputData, responses chan<- *web.Message) {
	for { // ever TODO close by control
		if value, ok := output.Receive(); !ok {
			responses <- web.ErrorMessage(fmt.Errorf("output %#v closed", data))
			return
		} else if bytes, err := json.Marshal(value); err != nil {
			responses <- web.ErrorMessage(fmt.Errorf("marshal output %#v failed: %v", data, err))
			return
		} else {
			responses <- marshal("output", web.OutputPayload{OutputData: data, Value: json.RawMessage(bytes)})
		}
	}
}

func dispatch(responses chan<- *web.Message, request *web.Message) *web.Message {
	switch request.Type {
	case "command":
		var command web.Command
		if err := json.Unmarshal(request.Payload, &command); err != nil {
			return web.ErrorMessage(fmt.Errorf("invalid command payload: %q", string(request.Payload)))
		}
		switch command.Call {
		case "flow":
			var data web.FlowData
			if err := json.Unmarshal(command.Data, &data); err != nil {
				return web.ErrorMessage(fmt.Errorf("invalid flow data: %q", string(command.Data)))
			}
			return web.FlowMessage(api.Instance(data.Name))
		case "add":
			var data web.AddData
			if err := json.Unmarshal(command.Data, &data); err != nil {
				return web.ErrorMessage(fmt.Errorf("invalid add call: %q", string(command.Data)))
			}
			if err := api.Instance(data.Flow).Add(data.Name, data.Component); err != nil {
				return web.ErrorMessage(fmt.Errorf("add call failed: %v", err))
			}
			return request
		case "connect":
			var data web.ConnectData
			if err := json.Unmarshal(command.Data, &data); err != nil {
				return web.ErrorMessage(fmt.Errorf("invalid connect call: %q", string(command.Data)))
			}
			if err := api.Instance(data.Flow).Connect(data.Out.Name, data.Out.Port, data.In.Name, data.In.Port); err != nil {
				return web.ErrorMessage(fmt.Errorf("connect call failed: %v", err))
			}
			return request
		case "input":
			var data web.InputData
			var value interface{}
			if err := json.Unmarshal(command.Data, &data); err != nil {
				return web.ErrorMessage(fmt.Errorf("invalid input call: %q", string(command.Data)))
			} else if err := json.Unmarshal(data.Value, &value); err != nil {
				return web.ErrorMessage(fmt.Errorf("value unmarshal failed: %v", err))
			} else if input, err := api.Instance(data.Flow).In(data.Name, data.Port); err != nil {
				return web.ErrorMessage(err)
			} else {
				go input.Send(value)
			}
			return request
		case "output":
			var data web.OutputData
			if err := json.Unmarshal(command.Data, &data); err != nil {
				return web.ErrorMessage(fmt.Errorf("invalid output call: %q", string(command.Data)))
			} else if output, err := api.Instance(data.Flow).Out(data.Name, data.Port); err != nil {
				return web.ErrorMessage(err)
			} else {
				go forward(output, data, responses)
			}
			return request
		case "flows":
			f := map[string]string{}
			for _, name := range api.Flows() {
				f[name] = "desc " + " " + name
			}
			return marshal("flows", f)
		case "components":
			c := map[string]string{}
			for _, name := range api.Components() {
				c[name] = "desc " + " " + name
			}
			return marshal("components", c)
		default:
			return web.ErrorMessage(fmt.Errorf("invalid command call: %q", command.Call))
		}
	default:
		return web.ErrorMessage(fmt.Errorf("unknown dispatch: %q", string(request.Payload)))
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
	requests := make(chan *web.Message)
	defer close(requests)
	out := make(chan struct{})
	defer close(out)
	go func() {
		for {
			var request web.Message
			if err := conn.ReadJSON(&request); err != nil {
				log.Println("read failed:", err)
				out <- struct{}{}
				return
			}
			requests <- &request
		}
	}()

	responses := make(chan *web.Message)
	defer close(responses)
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()
	for {
		select {
		case <-out:
			return
		case response := <-responses:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(&response); err != nil {
				log.Println("write failed:", err)
				return
			}
		case request := <-requests:
			response := dispatch(responses, request)
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
