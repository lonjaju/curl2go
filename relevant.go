package curl2go

import (
	"fmt"
	"strings"
)

type Proxy struct {
	Server string
	Auth   *BasicAuth
}

// Relevant is some of the parameters related to constructing a Go request
type Relevant struct {
	URL       string
	Method    string
	Headers   map[string]string
	Data      RelevantData
	DataType  string
	Insecure  bool
	BasicAuth *BasicAuth

	Host  string
	Proxy *Proxy
}

// ExtractRelevant parse and extract the relevant parameters from map to a struct
func ExtractRelevant(flags ParsedFlags) (*Relevant, error) {
	relevant := Relevant{
		Headers: make(map[string]string, 0),
		Data: RelevantData{
			Files: []string{},
		},
	}

	if len(flags.UnFlags) < 2 {
		return nil, fmt.Errorf("not curl")
	}

	relevant.URL = flags.UnFlags[1]

	if ok := flags.BoolFlags["head"]; ok {
		relevant.Method = "HEAD"
	}

	if headers, ok := flags.StringsFlags["header"]; ok {
		relevant.Headers = ParseHeaders(headers)
	}

	cmdRequest := flags.StringsFlags["request"]
	cmdDataBinary := flags.StringsFlags["data-binary"]
	cmdDataRaw := flags.StringsFlags["data-raw"]

	isRawPost := false

	if vals := cmdRequest; len(vals) > 0 {
		// if multiple, use last (according to curl docs
		relevant.Method = vals[len(vals)-1]
	} else if len(cmdDataBinary) > 0 {
		isRawPost = true
	} else if len(cmdDataRaw) > 0 {
		isRawPost = true
	}

	if isRawPost {
		// for --data-binary and --data-raw, use method POST & data-type raw
		relevant.Method = "POST"
		relevant.DataType = "raw"
	}

	dataAscii := make([]string, 0)
	dataFiles := make([]string, 0)

	loadData := func(d []string, dataRawFlag bool) {
		if relevant.Method == "" {
			relevant.Method = "POST"
		}

		if _, ok := relevant.Headers["Content-Type"]; !ok {
			relevant.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}

		for _, value := range d {
			if len(value) > 0 && value[0] == '@' && !dataRawFlag {
				dataFiles = append(dataFiles, value[1:])
			} else {
				dataAscii = append(dataAscii, value)
			}
		}
	}

	if cmdData, ok := flags.StringsFlags["data"]; ok {
		loadData(cmdData, false)
	}
	if cmdDataBinary, ok := flags.StringsFlags["data-binary"]; ok {
		loadData(cmdDataBinary, false)
	}
	if cmdDataRaw, ok := flags.StringsFlags["data-raw"]; ok {
		loadData(cmdDataRaw, true)
	}

	relevant.Data.Ascii = strings.Join(dataAscii, "&")
	relevant.Data.Files = dataFiles

	if cmdUser, ok := flags.StringsFlags["user"]; ok && len(cmdUser) > 0 {
		basicAuthString := cmdUser[len(cmdUser)-1]
		if basicAuthSplit := strings.Index(basicAuthString, ":"); basicAuthSplit > -1 {
			relevant.BasicAuth = &BasicAuth{
				User: basicAuthString[:basicAuthSplit],
				Pass: basicAuthString[basicAuthSplit+1:],
			}
		} else {
			relevant.BasicAuth = &BasicAuth{
				User: basicAuthString,
				Pass: "",
			}
		}
	}

	if cmdProxy, ok := flags.StringsFlags["proxy"]; ok && len(cmdProxy) > 0 {
		cmdProxyUser := flags.StringsFlags["proxy-user"]

		var basicAuth BasicAuth
		if len(cmdProxyUser) > 0 {
			basicAuthString := cmdProxy[len(cmdProxyUser)-1]

			if basicAuthSplit := strings.Index(basicAuthString, ":"); basicAuthSplit > -1 {
				basicAuth.User = basicAuthString[:basicAuthSplit]
				basicAuth.Pass = basicAuthString[basicAuthSplit+1:]
			} else {
				basicAuth.User = basicAuthString
			}
		}

		relevant.Proxy = &Proxy{
			Server: cmdProxy[len(cmdProxy)-1],
			Auth:   &basicAuth,
		}
	}

	if relevant.Method == "" {
		relevant.Method = "GET"
	}

	if _, ok := flags.BoolFlags["insecure"]; ok {
		relevant.Insecure = true
	}

	return &relevant, nil
}
