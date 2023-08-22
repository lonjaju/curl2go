package curl2go

import "testing"

func TestExtractRelevant(t *testing.T) {
	flags := ParsedFlags{
		UnFlags: []string{"curl", "https://example.com"},
		BoolFlags: map[string]bool{
			"head":     true,
			"insecure": true,
		},
		StringsFlags: map[string][]string{
			"header":     {"Content-Type: application/json"},
			"request":    {"POST"},
			"data":       {"key1=value1", "key2=value2"},
			"user":       {"username:password"},
			"proxy":      {"http://proxyserver.com"},
			"proxy-user": {"proxy_username:proxy_password"},
		},
	}

	relevant, err := ExtractRelevant(flags)
	if err != nil {
		t.Fatal(err)
	}

	expected := &Relevant{
		URL:     "https://example.com",
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		Data: RelevantData{
			Ascii: "key1=value1&key2=value2",
		},
		BasicAuth: &BasicAuth{
			User: "username",
			Pass: "password",
		},
		Proxy: &Proxy{
			Server: "http://proxyserver.com",
			Auth: &BasicAuth{
				User: "proxy_username",
				Pass: "proxy_password",
			},
		},
		Insecure: true,
	}

	if relevant.URL != expected.URL {
		t.Errorf("ExtractRelevant URL returned %+v, expected %+v", relevant.URL, expected.URL)
	}
	if relevant.Method != expected.Method {
		t.Errorf("ExtractRelevant Method returned %+v, expected %+v", relevant.Method, expected.Method)
	}

	if relevant.Headers["Content-Type"] != expected.Headers["Content-Type"] {
		t.Errorf("ExtractRelevant Headers returned %+v, expected %+v", relevant.Headers, expected.Headers)
	}

	if relevant.Data.Ascii != expected.Data.Ascii {
		t.Errorf("ExtractRelevant Data returned %+v, expected %+v", relevant.Data.Ascii, expected.Data.Ascii)
	}

	if relevant.BasicAuth.User != expected.BasicAuth.User {
		t.Errorf("ExtractRelevant BasicAuth.User returned %+v, expected %+v", relevant.BasicAuth.User, expected.BasicAuth.User)
	}
	if relevant.BasicAuth.Pass != expected.BasicAuth.Pass {
		t.Errorf("ExtractRelevant BasicAuth.Pass returned %+v, expected %+v", relevant.BasicAuth.Pass, expected.BasicAuth.Pass)
	}
	if relevant.Proxy.Server != expected.Proxy.Server {
		t.Errorf("ExtractRelevant Proxy.Server returned %+v, expected %+v", relevant.Proxy.Server, expected.Proxy.Server)
	}
	if relevant.Proxy.Auth.User != expected.Proxy.Auth.User {
		t.Errorf("ExtractRelevant Proxy.Auth.User returned %+v, expected %+v", relevant.Proxy.Auth.User, expected.Proxy.Auth.User)
	}
	if relevant.Proxy.Auth.Pass != expected.Proxy.Auth.Pass {
		t.Errorf("ExtractRelevant Proxy.Auth.Pass returned %+v, expected %+v", relevant.Proxy.Auth.Pass, expected.Proxy.Auth.Pass)
	}
	if relevant.Insecure != expected.Insecure {
		t.Errorf("ExtractRelevant Insecure returned %+v, expected %+v", relevant.Insecure, expected.Insecure)
	}
}
