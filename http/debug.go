package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"moul.io/http2curl"
)

// SendRequestDebug log CURL command and related response for debugging
// friendly.
func SendRequestDebug(req *http.Request) (raw []byte, httpResponse *http.Response, err error) {
	curlCommand, err := http2curl.GetCurlCommand(req)
	if err != nil {
		fmt.Printf("[imbot/http][%v] GetCurlCommand error: %v\n", time.Now(), err)
		return
	}

	curlLiteral := curlCommand.String()

	fmt.Printf("[imbot/http][%v] request CURL: %v\n", time.Now(), curlLiteral)

	// send
	httpResponse, err = (&http.Client{}).Do(req)
	if err != nil {
		fmt.Printf("[imbot/http][%v] request error: %v. CURL: %v\n", time.Now(), err, curlLiteral)
		return
	}
	defer httpResponse.Body.Close()

	raw, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		fmt.Printf("[imbot/http][%v] read body error: %v. CURL: %v\n", time.Now(), err, curlLiteral)
		return
	}

	fmt.Printf("[%v] response body: %s. CURL: %v\n", time.Now(), raw, curlLiteral)
	return
}
