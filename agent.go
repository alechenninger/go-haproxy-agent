package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"

	"github.com/open-policy-agent/opa/rego"
)

var r = rego.New(
	rego.Query("x = data.example.allow"),
	rego.Load([]string{"./example.rego"}, nil))

var ctx = context.TODO()

func main() {
	log.Print("listen 9000")

	listener, err := net.Listen("tcp4", "127.0.0.1:9000")
	if err != nil {
		log.Printf("error create listener, %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	handler := func(req *request.Request) {

		log.Printf("handle request EngineID: '%s', StreamID: '%d', FrameID: '%d' with %d messages\n", req.EngineID, req.StreamID, req.FrameID, req.Messages.Len())

		// messageName := "get-ip-reputation"

		// mes, err := req.Messages.GetByName(messageName)
		// if err != nil {
		// 	log.Printf("message %s not found: %v", messageName, err)
		// 	return
		// }

		// mes.KV.Get("test")

		// ipValue, ok := mes.KV.Get("ip")
		// if !ok {
		// 	log.Printf("var 'ip' not found in message")
		// 	return
		// }

		// ip, ok := ipValue.(net.IP)
		// if !ok {
		// 	log.Printf("var 'ip' has wrong type. expect IP addr")
		// 	return
		// }

		// ipScore := rand.Intn(100)

		// log.Printf("IP: %s, send score '%d'", ip.String(), ipScore)
		rs, err := r.Eval(ctx)
		if err != nil {
			log.Printf("Error evaluating rego %v", err)
			return
		}

		req.Actions.SetVar(action.ScopeSession, "test", rs[0].Bindings["x"])
	}

	a := agent.New(handler)

	if err := a.Serve(listener); err != nil {
		log.Printf("error agent serve: %+v\n", err)
	}
}
