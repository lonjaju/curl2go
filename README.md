# curl2go
Convert curl to golang code

Usage:
```shell
echo "curl -X POST https://reqbin.com/echo/post/json                                                                                  main ✚ ◼
-H "Content-Type: application/json"
-d '{"productId": 123456, "quantity": 100}' " | ./curl2go

```

input:
```shell
curl -X POST https://reqbin.com/echo/post/json                                                                                  main ✚ ◼
        -H "Content-Type: application/json"
        -d '{"productId": 123456, "quantity": 100}'
```

output:
```go
package main

import (
        "fmt"
        "log"
        "time"
        "net/http"
        "net/http/httputil"
)

func init() {
    log.SetFlags(log.Lshortfile | log.Ltime)
}

func main() {
        client := http.Client{
                Transport: &http.Transport{
                        Proxy:                 nil,
                        TLSHandshakeTimeout:   time.Second * 5,
                        DisableKeepAlives:     true,
                        ResponseHeaderTimeout: time.Second * time.Duration(10),
                },
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                        return http.ErrUseLastResponse
                },
                Timeout: time.Second * time.Duration(10),
        }

        
        r, _ := http.NewRequest("POST", "https://reqbin.com/echo/post/json", nil)
        

        resp, err := client.Do(r)
        if err != nil {
                log.Fatal("SendRequest", err.Error())
        }
        
        defer resp.Body.Close()
        data, err := httputil.DumpResponse(resp, true)
        if err != nil {
                log.Fatal("DumpResponse", err)
        }

        fmt.Println(string(data))
}
```

```