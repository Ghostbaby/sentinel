package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type BaseClient struct {
	HTTP      *http.Client
	Endpoint  string
	Transport *http.Transport
}

func NewBaseClient(url string, timeout time.Duration) *BaseClient {
	hClient := &http.Client{}
	if timeout > 0 {
		hClient.Timeout = timeout
	}

	client := &BaseClient{
		HTTP:      hClient,
		Endpoint:  url,
		Transport: &http.Transport{},
	}

	return client
}

func (c *BaseClient) Get(ctx context.Context, pathWithQuery string, out interface{}) error {
	return c.request(ctx, http.MethodGet, pathWithQuery, nil, out)
}

func (c *BaseClient) Post(ctx context.Context, pathWithQuery string, in, out interface{}) error {
	return c.request(ctx, http.MethodPost, pathWithQuery, in, out)
}

func (c *BaseClient) Delete(ctx context.Context, pathWithQuery string, in interface{}) error {
	return c.request(ctx, http.MethodDelete, pathWithQuery, in, nil)
}

func (c *BaseClient) request(
	ctx context.Context,
	method string,
	pathWithQuery string,
	requestObj,
	responseObj interface{},
) error {
	var body io.Reader = http.NoBody
	if requestObj != nil {
		outData, err := json.Marshal(requestObj)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(outData)
	}

	request, err := http.NewRequest(method, Joins(c.Endpoint, pathWithQuery), body)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if responseObj != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseObj); err != nil {
			return err
		}
	}

	return nil
}

func (c *BaseClient) doRequest(context context.Context, request *http.Request) (*http.Response, error) {
	withContext := request.WithContext(context)

	response, err := c.HTTP.Do(withContext)
	if err != nil {
		fmt.Println(err)
		return response, err
	}

	err = checkError(response)
	return response, err
}

func checkError(response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		data, _ := ioutil.ReadAll(response.Body)
		return errors.New(
			//fmt.Sprintf("请求执行失败, 返回码 %d error: %s", response.StatusCode, string(data)))
			string(data))
	}
	return nil
}

func (c *BaseClient) Close() {
	if c.Transport != nil {
		// When the http transport goes out of scope, the underlying goroutines responsible
		// for handling keep-alive connections are not closed automatically.
		// Since this client gets recreated frequently we would effectively be leaking goroutines.
		// Let's make sure this does not happen by closing idle connections.
		c.Transport.CloseIdleConnections()
	}
}

func (c *BaseClient) Equal(c2 *BaseClient) bool {
	// handle nil case
	if c2 == nil && c != nil {
		return false
	}

	// compare endpoint and user creds
	return c.Endpoint == c2.Endpoint
}

func Joins(args ...string) string {
	var str strings.Builder
	for _, arg := range args {
		str.WriteString(arg)
	}
	return str.String()
}
