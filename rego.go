package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/request"

	"github.com/open-policy-agent/opa/rego"
)

func regoHandler(rego *rego.Rego, opts ...rego.PrepareOption) func(context.Context, *request.Request) {
	query, err := rego.PrepareForEval(context.Background(), opts...)
	if err != nil {
		// TODO: don't log fatal
		log.Fatal(err)
	}

	return regoHandlerForQuery(query)
}

// TODO: test this
func regoHandlerForQuery(query rego.PreparedEvalQuery) func(context.Context, *request.Request) {
	return func(ctx context.Context, req *request.Request) {
		// TODO: some of this could be shared for other kinds of handlers (that didn't use OPA)
		log.Printf("handle request EngineID: '%s', StreamID: '%d', FrameID: '%d' with %d messages\n",
			req.EngineID, req.StreamID, req.FrameID, req.Messages.Len())

		mes, err := req.Messages.GetByName("goagent")
		if err != nil {
			log.Printf("no goagent message: %v", err)
			return
		}

		args := mes.KV.Data()
		input := make(map[string]interface{}, len(args))

		for _, arg := range args {
			// TODO: configurable header argument name?
			if arg.Name == "header" {
				hdrBytes := arg.Value.([]byte)
				input[arg.Name], _, err = headers(hdrBytes)
				if err != nil {
					log.Printf("Error parsing headers %v", err)
					return
				}
			} else if arg.Name == "body" {
				// TODO: configurable body parsing?
				// Look at content type header?
				var body interface{}
				json.Unmarshal(arg.Value.([]byte), &body)
				input[arg.Name] = body
			} else {
				input[arg.Name] = arg.Value
			}
		}

		log.Println(input)

		rs, err := query.Eval(ctx, rego.EvalInput(input))
		if err != nil {
			log.Printf("Error evaluating rego %v", err)
			return
		}

		results := len(rs)
		if results == 0 {
			return
		}

		for k, v := range rs[0].Bindings {
			log.Printf("%v : %v", k, v)
			req.Actions.SetVar(action.ScopeRequest, k, v)
		}
	}
}
