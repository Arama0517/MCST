package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Request 请求URL, 返回响应; 运行成功后请添加`defer resp.Body.Close()`到你的代码内
func Request(URL url.URL, Method string, Header http.Header) (*http.Response, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(Method, URL.String(), nil)
	req.Header = Header
	req.Header.Set("User-Agent", fmt.Sprintf("MCSCS-Go/%s", VERSION))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}
	return resp, nil
}
