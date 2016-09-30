package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"gopkg.in/elazarl/goproxy.v1"
	"gopkg.in/elazarl/goproxy.v1/transport"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var proxyHeader = "X-proxy-req-id"

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	keyParam := flag.String("key", "X-tid", "which header attribute should be used as request key")
	proxyAddr := flag.String("proxy", "8080", "proxy listen address")
	serverAddr := flag.String("server", "8081", "server listen address")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	tr := transport.Transport{Proxy: transport.ProxyFromEnvironment}

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
			ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
			return
		})
		id := fmt.Sprintf("%v-%v", time.Now().UnixNano(), ctx.Session)
		uri := req.URL.Host + req.URL.Path
		keyValue := req.Header.Get(*keyParam)
		Info.Printf("saving request uri:[%s] key:[%s] id:[%s]", uri, keyValue, id)
		Set(createMainCacheId(keyValue, uri, strings.ToUpper(req.Method)), id)
		req.Header.Add(proxyHeader, id)

		b, body := readBody(req.Body)
		req.Body = body
		content := HttpContent{ID: id, Timestamp: time.Now().UnixNano() / int64(time.Millisecond), URI: req.RequestURI, Method: req.Method, Body: b, Headers: toHeaders(req.Header)}
		SetHttpContent("req-"+id, content)
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil {
			return resp
		}

		b, body := readBody(resp.Body)
		resp.Body = body
		id := resp.Request.Header.Get(proxyHeader)
		Info.Printf("saving response id:[%s] code:[%v]", id, resp.StatusCode)
		content := HttpContent{ID: id, Timestamp: time.Now().UnixNano() / int64(time.Millisecond), Code: resp.StatusCode, Body: b, Headers: toHeaders(resp.Header)}
		SetHttpContent("resp-"+id, content)
		return resp
	})

	go func(addr string) {
		e := echo.New()
		e.SetHTTPErrorHandler(errorHandler)
		e.Get("/version", version)
		e.Get("/metadata/:type", metadata())
		e.Get("/body/:type", body)
		Info.Printf("listening server on %v", addr)
		server := standard.New(fmt.Sprintf(":%v", addr))
		server.SetHandler(e)
		server.ListenAndServe()
	}(*serverAddr)

	Info.Printf("listening proxy on %v", *proxyAddr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *proxyAddr), proxy))
}

func createMainCacheId(tid, uri, method string) string {
	return "id-" + method + "-" + tid + "-" + uri
}

func version(c echo.Context) error {
	return c.String(http.StatusOK, "v1.0, "+runtime.Version())
}

func metadata() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := c.Param("type")
		if "req" != param && "resp" != param {
			return fmt.Errorf("type:[%v] not allowed", param)
		}

		key := c.QueryParam("key")
		uri := c.QueryParam("uri")
		method := strings.ToUpper(c.QueryParam("method"))

		Info.Printf("[%s] get meta info for request type:[%s] method:[%s] key:[%s] uri:[%s]", c.Request().URL().Path(), param, method, key, uri)

		id, err := Get(createMainCacheId(key, uri, method))
		if err != nil {
			return err
		}
		content, err := GetHttpContent(param + "-" + id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, supressBody(content))
	}

}

func supressBody(content HttpContent) HttpContent {
	content.Body = nil
	return content
}

func body(c echo.Context) error {
	param := c.Param("type")
	if "req" != param && "resp" != param {
		return fmt.Errorf("type:[%v] not allowed", param)
	}

	tid := c.QueryParam("id")
	uri := c.QueryParam("uri")
	method := strings.ToUpper(c.QueryParam("method"))

	Info.Printf("[%s] get body for request type:[%s] method:[%s] tid:[%s] uri:[%s]", c.Request().URI(), param, method, tid, uri)
	id, err := Get(createMainCacheId(tid, uri, method))
	if err != nil {
		return err
	}
	content, err := GetHttpContent(param + "-" + id)
	if err != nil {
		return err
	}
	for k, v := range content.Headers {
		c.Response().Header().Add(k, v)
	}
	if content.Code == 0 {
		c.Response().WriteHeader(http.StatusOK)
	} else {
		c.Response().WriteHeader(content.Code)
	}
	c.Response().Write(content.Body)
	return nil
}

func errorHandler(err error, c echo.Context) {
	switch err := err.(type) {

	}
}

func readBody(body io.ReadCloser) ([]byte, io.ReadCloser) {
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}
	return reqBody, WrapperReader{bytes.NewBuffer(reqBody)}
}

func toHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, value := range headers {
		result[key] = value[0]
	}
	return result
}

type WrapperReader struct {
	*bytes.Buffer
}

func (WrapperReader) Close() error { return nil }

type HttpContent struct {
	ID        string            `json:"id,omitempty"`
	Timestamp int64             `json:"timestamp,omitempty"`
	URI       string            `json:"uri,omitempty"`
	Method    string            `json:"method,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Code      int               `json:"code,omitempty"`
	Body      []byte            `json:"body,omitempty"`
}
