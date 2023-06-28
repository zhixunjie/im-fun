package http

import (
	"net"
	"net/http"
	"time"
)

type Options struct {
	// TCP
	tcpTimeout   time.Duration
	tcpKeepAlive time.Duration

	// connection
	idleConnTimeout     time.Duration
	maxIdleConns        int
	maxIdleConnsPerHost int

	// request
	requestTimeout time.Duration
}

type Option func(*Options)

func TcpTimeout(tcpTimeout time.Duration) Option {
	return func(o *Options) {
		o.tcpTimeout = tcpTimeout
	}
}

func TcpKeepAlive(tcpKeepAlive time.Duration) Option {
	return func(o *Options) {
		o.tcpKeepAlive = tcpKeepAlive
	}
}
func IdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleConnTimeout = idleConnTimeout
	}
}
func MaxIdleConns(maxIdleConns int) Option {
	return func(o *Options) {
		o.maxIdleConns = maxIdleConns
	}
}
func MaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(o *Options) {
		o.maxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func RequestTimeout(requestTimeout time.Duration) Option {
	return func(o *Options) {
		o.requestTimeout = requestTimeout
	}
}

// NewClient NewClient NewClient NewClient NewClient

func NewClient(optionFunc ...Option) *Client {
	client := new(Client)

	// 创建option，并设定默认值
	newOptions := Options{
		tcpTimeout:          30 * time.Second,
		tcpKeepAlive:        30 * time.Second,
		idleConnTimeout:     90 * time.Second,
		maxIdleConns:        100,
		maxIdleConnsPerHost: http.DefaultMaxIdleConnsPerHost,
		requestTimeout:      time.Second * 5,
	}

	// 根据传入的回调函数切片，调用回调函数，覆盖默认值
	for _, o := range optionFunc {
		o(&newOptions)
	}

	// create object: http.Client
	// refer: /Users/jasonzhi/sdk/go1.18.10/src/net/http/transport.go:43
	client.httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// TCP: 控制TCP连接相关
			DialContext: (&net.Dialer{
				Timeout:   newOptions.tcpTimeout,   // TCP建立连接的超时时间
				KeepAlive: newOptions.tcpKeepAlive, // 设置了活跃连接的TCP-KeepAlive探针间隔
			}).DialContext,
			ForceAttemptHTTP2: true,
			// 空闲连接: 控制连接池子中的空闲连接
			IdleConnTimeout:     newOptions.tcpKeepAlive,        // 空闲连接KeepAlive的超时时间（超时后回自动断开）
			MaxIdleConns:        newOptions.maxIdleConns,        // 最大的空闲连接
			MaxIdleConnsPerHost: newOptions.maxIdleConnsPerHost, // 最大的空闲连接(每个host)，这个参数默认为2，一般情况下不需要设置
		},
		Timeout: newOptions.requestTimeout, // 一个HTTP请求过程的超时时间
	}

	// set other property
	client.needBody = true

	return client
}

/**
调研：
http.Client、http.Transport的具体研究
- https://duyanghao.github.io/http-transport/
- https://www.cnblogs.com/charlieroro/p/11409153.html

q1: Transport.IdleConnTimeout与net.Dialer.KeepAlive有什么关系，哪一个是所谓的HTTP keep-alives？？？
*/
func newClient() *Client {
	client := new(Client)

	client.httpClient = &http.Client{
		Transport: &http.Transport{
			// TCP: 控制TCP连接相关
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,                 // TCP建立连接的超时时间
				Deadline:  time.Now().Add(30 * time.Second), // TCP建立连接的超时时间(跟参数Timeout一样的效果，只是值的方式不一样) 注意：配置了这个值后，HTTP请求一段时间后会报错：dial timeout
				KeepAlive: 30 * time.Second,                 // 设置了活跃连接的TCP-KeepAlive探针间隔
			}).DialContext,
			// 空闲连接: 控制连接池子中的空闲连接
			IdleConnTimeout:     90 * time.Second, // 空闲连接KeepAlive的超时时间（超时后回自动断开）
			MaxIdleConns:        100,              // 最大的空闲连接
			MaxIdleConnsPerHost: 100,              // 最大的空闲连接(每个host)，这个参数默认为2，一般情况下不需要设置
		},
		Timeout: time.Second * 5, // 一个HTTP请求过程的超时时间
	}
	return client
}
