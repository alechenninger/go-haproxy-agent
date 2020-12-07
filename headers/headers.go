package headers

import (
	"net/http"

	"github.com/negasus/haproxy-spoe-go/typeddata"
	"github.com/negasus/haproxy-spoe-go/varint"
)

// ParseHeaders parses haproxy req.hdrs_bin into http.Header and number of bytes read out of buffer.
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
func ParseHeaders(buf []byte) (header http.Header, bytes int, err error) {
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
