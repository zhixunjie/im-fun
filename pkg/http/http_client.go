package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrorRequestNotSet = errors.New("please set request object")
var ErrorStatusNotOK = errors.New("status code is not ok")

type Client struct {
	httpClient *http.Client
	request    *http.Request

	needBody bool
}

type Response struct {
	Header   map[string][]string
	NeedBody bool
	Body     []byte
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

func (client *Client) SetNeedBody(need bool) {
	client.needBody = need
}

// Get 使用Client对象进行请求-GET
func (client *Client) Get(reqUrl string, params map[string]string, headers map[string]string) (rsp Response, err error) {
	logHead := "Get|"
	rsp = Response{
		NeedBody: client.needBody,
	}

	// get consume time
	start := time.Now()
	defer func() {
		logging.Infof(logHead+"consume=%vms,reqUrl=%v", time.Now().Sub(start).Milliseconds(), reqUrl)
	}()

	// new request
	request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		logging.Errorf(logHead+"http.NewRequest err=%v\n", err)
		return rsp, err
	}
	client.setRequest(request)

	// add sth
	_ = client.addQueryParams(params) // add query params
	_ = client.addHeaders(headers)    // add headers

	// send request
	return client.do()
}

func (client *Client) PostJson(reqUrl string, body any, params map[string]string, headers map[string]string) (rsp Response, err error) {
	logHead := "PostJson|"
	rsp = Response{
		NeedBody: client.needBody,
	}

	// marshal body
	var bodyJson []byte
	if body != nil {
		bodyJson, err = json.Marshal(body)
		if err != nil {
			fmt.Printf(logHead+"json.Marshal err=%v\n", err)
			return rsp, err
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
func (client *Client) Post(reqUrl string, body io.Reader, params map[string]string, headers map[string]string) (rsp Response, err error) {
	logHead := "Post|"
	rsp = Response{
		NeedBody: client.needBody,
	}

	var bodyLog []byte
	if v, ok := body.(*bytes.Buffer); ok {
		bodyLog = v.Bytes()
	}

	// get consume time
	start := time.Now()
	defer func() {
		logging.Infof(logHead+"consume=%vms,reqUrl=%v,body=%s", time.Now().Sub(start).Milliseconds(), reqUrl, bodyLog)
	}()

	// new request
	request, err := http.NewRequest(http.MethodPost, reqUrl, body)
	if err != nil {
		fmt.Printf(logHead+"http.NewRequest err=%v\n", err)
		return rsp, err
	}
	client.setRequest(request)

	// add sth
	_ = client.addQueryParams(params) // add query params
	_ = client.addHeaders(headers)    // add headers

	// send request
	return client.do()
}

func (client *Client) do() (rsp Response, err error) {
	logHead := "do|"
	httpClient := client.httpClient
	request := client.request
	rsp = Response{
		NeedBody: client.needBody,
	}

	// log
	//logging.Infof(logHead+"method=%s,url=%s\n", request.Method, request.URL.String())

	// do request
	result, err := httpClient.Do(request)
	if err != nil {
		logging.Errorf(logHead+"client.Do err=%v\n", err)
		return
	}

	defer func() {
		_ = result.Body.Close()
	}()

	// get header
	rsp.Header = result.Header
	//PrintHeaderMap(result.Header)

	// check status code
	if result.StatusCode != http.StatusOK {
		err = ErrorStatusNotOK
		logging.Errorf(logHead+"status not allow, result.Status=%v", result.Status)
		return
	}

	// read content
	if client.needBody {
		body, newErr := ioutil.ReadAll(result.Body)
		if newErr != nil {
			err = newErr
			logging.Errorf(logHead+"ReadAll err=%v\n", newErr)
			return
		}
		rsp.Body = body
	}

	return
}
