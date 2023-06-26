package http

import (
	"fmt"
	"testing"
)

//var reqUrl = "http://www.jason.com/APISearch/cgi-bin/index.php"
var reqUrl = "http://www.baidu.com"

func TestHttpClientGet(t *testing.T) {
	c := NewClient()
	params := map[string]string{
		"a": "1",
		"b": "2",
	}
	fmt.Println(c.Get(reqUrl, params, nil))

}

func TestHttpClientPost(t *testing.T) {
	c := NewClient()
	params := map[string]string{
		"a": "1",
		"b": "2",
	}
	body := map[string]string{
		"b1": "100",
		"b2": "101",
	}

	fmt.Println(c.PostJson(reqUrl, body, params, nil))
}
