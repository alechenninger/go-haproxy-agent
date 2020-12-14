package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"
	"github.com/negasus/haproxy-spoe-go/typeddata"
	"github.com/negasus/haproxy-spoe-go/varint"
	"github.com/spf13/cobra"

	"github.com/open-policy-agent/opa/rego"
)

var rootCmd = &cobra.Command{
	Use:   "go-haproxy-agent",
	Short: "haproxy spoa agent for golang extensions",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Print("listen 9000")

		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return err
		}
		load, err := cmd.Flags().GetStringArray("load")
		if err != nil {
			return err
		}

		handler := regoHandler(rego.New(
			rego.Query(query),
			rego.Load(load, nil)))

		a := agent.New(func(r *request.Request) {
			// TODO: make deadline per request configurable
			//   should be same as haproxy agent timeout e.g. timeout processing 5s
			handler(context.TODO(), r)
		})

		listener, err := net.Listen("tcp4", "127.0.0.1:9000")
		if err != nil {
			log.Printf("error create listener, %v", err)
			return err
		}
		defer listener.Close()

		if err := a.Serve(listener); err != nil {
			log.Printf("error agent serve: %+v\n", err)
			return err
		}

		return nil
	},
}

func main() {
	rootCmd.Flags().String("query", "", "Rego query. Results will be set as variables in "+
		"haproxy, using the same key-value pairs as in the query results. Only the first "+
		"result in result set is used.")
	rootCmd.MarkFlagRequired("query")

	rootCmd.Flags().StringArray("load", []string{}, "Path to rego policy and input files to be "+
		"used in policy evaluation. May be passed more than once to load multiple paths.")
	rootCmd.MarkFlagRequired("load")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// header parses haproxy req.hdrs_bin into http.Header and number of bytes read out of buffer.
//
// From haproxy docs:
//
//     Each string is described by a length followed by the number of bytes indicated in the length.
//     The length is represented using the variable integer encoding detailed in the SPOE
//     documentation. The end of the list is marked by a couple of empty header names and values
//     (length of 0 for both).
//
//     *(<str:header-name><str:header-value>)<empty string><empty string>
//
//     int:  refer to the SPOE documentation for the encoding
//     str:  <int:length><bytes>
func header(buf []byte) (header http.Header, bytes int, err error) {
	header = make(http.Header)

	var name string
	var value string
	var b int

	for {
		name, b, err = str(buf)
		buf = buf[b:]
		bytes += b

		if err != nil {
			return nil, bytes, err
		}

		value, b, err = str(buf)
		buf = buf[b:]
		bytes += b

		if err != nil {
			return nil, bytes, err
		}

		if len(name) == 0 && len(value) == 0 {
			// Or return once len(buf) == 0?
			return
		}

		header.Add(name, value)
	}
}

// Parse string from buffer according to SPOP string type.
func str(buffer []byte) (str string, bytes int, err error) {
	length, bytes := varint.Uvarint(buffer)
	buffer = buffer[bytes:]
	if len(buffer) < int(length) {
		err = typeddata.ErrDecodingBufferTooSmall
		return
	}
	str = string(buffer[:length])
	bytes += int(length)
	return str, bytes, err
}
