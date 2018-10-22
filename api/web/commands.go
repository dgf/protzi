package web

import (
	"encoding/json"
	"fmt"
)

type Command struct {
	Call string          `json:"call"`
	Data json.RawMessage `json:"data"`
}

type FlowData struct {
	Name string `json:"name"`
}

type AddData struct {
	Flow      string `json:"flow"`
	Name      string `json:"name"`
	Component string `json:"component"`
}

type PortData struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type ConnectData struct {
	Flow string   `json:"flow"`
	In   PortData `json:"in"`
	Out  PortData `json:"out"`
}

type InputData struct {
	Flow  string          `json:"flow"`
	Name  string          `json:"name"`
	Port  string          `json:"port"`
	Value json.RawMessage `json:"value"`
}

type OutputData struct {
	Flow string `json:"flow"`
	Name string `json:"name"`
	Port string `json:"port"`
}

func ComponentsCommand() Command {
	return Command{Call: "components"}
}

func FlowsCommand() Command {
	return Command{Call: "flows"}
}

func FlowCommand(name string) Command {
	return Command{
		Call: "flow",
		Data: json.RawMessage([]byte(fmt.Sprintf(
			`{"name":%q}`, name))),
	}
}

func AddCommand(flow, name, component string) Command {
	return Command{
		Call: "add",
		Data: json.RawMessage([]byte(fmt.Sprintf(
			`{"flow":%q,"name":%q,"component":%q}`,
			flow, name, component))),
	}
}

func ConnectCommand(flow, outName, outPort, inName, inPort string) Command {
	return Command{
		Call: "connect",
		Data: json.RawMessage([]byte(fmt.Sprintf(
			`{"flow":%q,"out":{"name":%q,"port":%q},"in":{"name":%q,"port":%q}}`,
			flow, outName, outPort, inName, inPort))),
	}
}

func InputCommand(flow, name, port string, value json.RawMessage) Command {
	return Command{
		Call: "input",
		Data: json.RawMessage([]byte(fmt.Sprintf(
			`{"flow":%q,"name":%q,"port":%q,"value":%s}`,
			flow, name, port, value))),
	}
}

func OutputCommand(flow, name, port string) Command {
	return Command{
		Call: "output",
		Data: json.RawMessage([]byte(fmt.Sprintf(
			`{"flow":%q,"name":%q,"port":%q}`,
			flow, name, port))),
	}
}
