package lib

import (
	"fmt"
	"net/http"
	"net/url"
)

// Request 请求URL, 返回响应; 运行成功后请添加`defer resp.Body.Close()`到你的代码内
func Request(URL url.URL, Method string, Header map[string]string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest(Method, URL.String(), nil)
	req.Header.Set("User-Agent", fmt.Sprintf("MCSCS-Go/%s", VERSION))
	for k, v := range Header {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
