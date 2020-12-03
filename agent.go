package main

import (
	"context"
	"log"
	"net"
	"os"

	"alechenninger.com/haproxy-go-extensions/headers"

	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"

	"github.com/open-policy-agent/opa/rego"
)

var r = rego.New(
	rego.Query(os.Args[2]),
	rego.Load([]string{os.Args[1]}, nil))

var ctx = context.TODO()

func main() {
	log.Print("listen 9000")

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp4", "127.0.0.1:9000")
	if err != nil {
		log.Printf("error create listener, %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	// TODO: test this
	handler := func(req *request.Request) {
		log.Printf("handle request EngineID: '%s', StreamID: '%d', FrameID: '%d' with %d messages\n", req.EngineID, req.StreamID, req.FrameID, req.Messages.Len())

		mes, err := req.Messages.GetByName("goagent")
		if err != nil {
			log.Printf("no goagent message: %v", err)
			return
		}

		items := mes.KV.Data()
		input := make(map[string]interface{}, len(items))

		for _, i := range items {
			// TODO: configurable header argument name?
			if i.Name == "header" {
				hdrBytes := i.Value.([]byte)
				input[i.Name], _, err = headers.ParseHeaders(hdrBytes)
				if err != nil {
					log.Printf("Error parsing headers %v", err)
					return
				}
			} else {
				input[i.Name] = i.Value
			}
		}

		log.Println(input)

		rs, err := query.Eval(ctx, rego.EvalInput(input))
		if err != nil {
			log.Printf("Error evaluating rego %v", err)
			return
		}

		if len(rs) == 0 {
			return
		}

		for k, v := range rs[0].Bindings {
			req.Actions.SetVar(action.ScopeSession, k, v)
		}
	}

	a := agent.New(handler)

	if err := a.Serve(listener); err != nil {
		log.Printf("error agent serve: %+v\n", err)
	}
}