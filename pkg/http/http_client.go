package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var ErrorRequestNotSet = errors.New("please set request object")

type Client struct {
	httpClient *http.Client
	request    *http.Request
}

// NewClient 创建Client
func NewClient() *Client {
	client := new(Client)

	client.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // TCP建立连接的超时时间
				KeepAlive: 30 * time.Second, // 设置了活跃连接的TCP-KeepAlive探针间隔
				//Deadline:  time.Now().Add(30 * time.Second), // TCP建立连接的超时时间(跟参数Timeout一样的效果)
			}).DialContext,
			// 控制连接池子中，空闲连接的参数
			IdleConnTimeout: 90 * time.Second, // 空闲连接KeepAlive的超时时间（超时后回自动断开）
			MaxIdleConns:    100,              // 最大的空闲连接
		},
		Timeout: time.Second * 5, // 一个HTTP请求过程的超时时间
	}
	return client
}

func (client *Client) setRequest(request *http.Request) {
	client.request = request
}

// 添加请求头(header)
func (client *Client) addHeaders(headers map[string]string) error {
	request := client.request
	if request == nil {
		return ErrorRequestNotSet
	}
	// set value
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	return nil
}

// 添加查询参数(GET、POST)
func (client *Client) addQueryParams(params map[string]string) error {
	request := client.request
	if request == nil {
		return ErrorRequestNotSet
	}
	// set value
	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	return nil
}

// Get 使用Client对象进行请求-GET
func (client *Client) Get(reqUrl string, params map[string]string, headers map[string]string) error {
	logHead := "httpGet|"

	// new request
	request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		logging.Errorf(logHead+"http.NewRequest err=%v\n", err)
		return err
	}
	client.setRequest(request)

	// add sth
	_ = client.addQueryParams(params) // add query params
	_ = client.addHeaders(headers)    // add headers

	// send request
	return client.do()
}

func (client *Client) PostJson(reqUrl string, body any, params map[string]string, headers map[string]string) error {
	logHead := "PostJson|"
	var bodyJson []byte
	var err error

	// marshal body
	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			fmt.Printf(logHead+"json.Marshal err=%v\n", err)
			return err
		}
	}

	// set header
	if len(headers) == 0 {
		headers = make(map[string]string, 10)
	}
	headers["Content-type"] = "application/json"

	return client.Post(reqUrl, bytes.NewBuffer(bodyJson), params, headers)
}

// Post 使用Client对象进行请求-POST
func (client *Client) Post(reqUrl string, body io.Reader, params map[string]string, headers map[string]string) error {
	logHead := "Post|"
	var err error

	// new request
	request, err := http.NewRequest(http.MethodPost, reqUrl, body)
	if err != nil {
		fmt.Printf(logHead+"http.NewRequest err=%v\n", err)
		return err
	}
	client.setRequest(request)

	// add sth
	_ = client.addQueryParams(params) // add query params
	_ = client.addHeaders(headers)    // add headers

	// send request
	return client.do()
}

func (client *Client) do() error {
	logHead := "do|"
	httpClient := client.httpClient
	request := client.request
	logging.Infof(logHead+"method=%s,url=%s\n", request.Method, request.URL.String())

	// do request
	resp, err := httpClient.Do(request)
	if err != nil {
		logging.Errorf(logHead+"client.Do err=%v\n", err)
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// read content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Errorf(logHead+"ReadAll err=%v\n", err)
		return err
	}
	logging.Infof(logHead+"ReadAll body=%v\n", string(body))

	return nil
}
