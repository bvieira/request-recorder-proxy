Request Recorder Proxy
======================
tool using [go] to help automated tests validations for asynchronous services

# Flags:
Use '-h' for flags usage
```sh
request-recorder-proxy -h
```

# API:
- [Version](#version)
- [Requests](#requests)
- [Metadata](#metadata)
- [Body](#body)

## Version
returns application version

### Request:
`GET` /version

### Response:
	string

```sh
$ curl -v localhost:8081/version 
> GET /version HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Sat, 01 Oct 2016 00:00:07 GMT
< Content-Length: 11
<
v1.0, go1.7
```

## Requests
returns the list of ids for requests done with key, uri and method

### Request:
`GET` /requests?key=&uri=&method=

| query parameter   | description           |
|-------------------|-----------------------|
| `key`             | `key` header used on request  |
| `uri`             | uri without query parameters used on request  |
| `method`          | http method used on request    |


### Response:
	[]string
```sh
$ curl -v "localhost:8081/requests?key=test&uri=www.google.com/&method=get"
> GET /requests?key=test&uri=www.google.com/&method=get HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sat, 01 Oct 2016 00:10:39 GMT
< Content-Length: 73
<
["1475276740074193887-4","1475276741289460299-5","1475276742610169482-6"]
```

## Metadata
Return metadata info for :id and :type

### Request:
`GET` /metadata/:id/:type

| param   | description           |
|-------------------|-----------------------|
| `:id`             | `id` returned from /requests  |
| `:type`             | 'req' for request and 'resp' for response  |

### Response:
- [Metadata Request](#metadata-request)
	```sh
	$ curl -v "localhost:8081/metadata/1475276740074193887-4/req"
	> GET /metadata/1475276740074193887-4/req HTTP/1.1
	> Host: localhost:8081
	> User-Agent: curl/7.43.0
	> Accept: */*
	>
	< HTTP/1.1 200 OK
	< Content-Type: application/json; charset=utf-8
	< Date: Sat, 01 Oct 2016 00:22:12 GMT
	< Content-Length: 244
	<
	{"id":"1475276740074193887-4","timestamp":1475276740076,"uri":"http://www.google.com/","method":"GET","headers":{"Accept":"*/*","Proxy-Connection":"Keep-Alive","User-Agent":"curl/7.43.0","X-Proxy-Req-Id":"1475276740074193887-4","X-Tid":"test"}}
	```

- [Metadata Response](#metadata-response)
	```sh
	$ curl -v "localhost:8081/metadata/1475276740074193887-4/resp"
	> GET /metadata/1475276740074193887-4/resp HTTP/1.1
	> Host: localhost:8081
	> User-Agent: curl/7.43.0
	> Accept: */*
	>
	< HTTP/1.1 200 OK
	< Content-Type: application/json; charset=utf-8
	< Date: Sat, 01 Oct 2016 00:22:28 GMT
	< Content-Length: 289
	<
	{"id":"1475276740074193887-4","timestamp":1475276740100,"headers":{"Cache-Control":"private","Content-Length":"262","Content-Type":"text/html; charset=UTF-8","Date":"Fri, 30 Sep 2016 23:05:40 GMT","Location":"http://www.google.com.br/?gfe_rd=cr\u0026ei=xO_uV_-vBfPM8AeX-oG4Cg"},"code":302}
	```



## Body
Return what was sent and received on proxy for id

### Request:
`GET` /body/:id/:type

| param   | description           |
|-------------------|-----------------------|
| `:id`             | `id` returned from /requests  |
| `:type`             | 'req' for request and 'resp' for response  |

### Response:
- 'req'
	```sh
	$ curl -v "localhost:8081/body/1475276740074193887-4/req"
	> GET /body/1475276740074193887-4/req HTTP/1.1
	> Host: localhost:8081
	> User-Agent: curl/7.43.0
	> Accept: */*
	>
	< HTTP/1.1 200 OK
	< Accept: */*
	< Proxy-Connection: Keep-Alive
	< User-Agent: curl/7.43.0
	< X-Proxy-Req-Id: 1475276740074193887-4
	< X-Tid: test
	< Date: Sat, 01 Oct 2016 00:33:38 GMT
	< Content-Length: 0
	< Content-Type: text/plain; charset=utf-8
	<
	```
- 'resp'
	```sh
	$ curl -v "localhost:8081/body/1475276740074193887-4/resp"
	> GET /body/1475276740074193887-4/resp HTTP/1.1
	> Host: localhost:8081
	> User-Agent: curl/7.43.0
	> Accept: */*
	>
	< HTTP/1.1 302 Found
	< Cache-Control: private
	< Content-Length: 262
	< Content-Type: text/html; charset=UTF-8
	< Date: Fri, 30 Sep 2016 23:05:40 GMT
	< Location: http://www.google.com.br/?gfe_rd=cr&ei=xO_uV_-vBfPM8AeX-oG4Cg
	<
	<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
	<TITLE>302 Moved</TITLE></HEAD><BODY>
	<H1>302 Moved</H1>
	The document has moved
	<A HREF="http://www.google.com.br/?gfe_rd=cr&amp;ei=xO_uV_-vBfPM8AeX-oG4Cg">here</A>.
	</BODY></HTML>
```


# Schema
## Metadata Request
```json
{
	"id": string,
	"timestamp": integer,
	"uri": string,
	"method": string,
	"headers": map[string]string
}
```

## Metadata Response
```json
{
	"id": string,
	"timestamp": integer,
	"headers": map[string]string,
	"code": integer
}
```

# Example:

Request using proxy
```sh
curl "http://www.google.com/" --proxy localhost:8080 -H "X-tid: test"
```

Get request id
```sh
curl -v "localhost:8081/requests?key=test&uri=www.google.com/&method=get"
```

Retrieve request/response infos
```sh
curl -v "localhost:8081/metadata/1475276740074193887-4/req"
curl -v "localhost:8081/metadata/1475276740074193887-4/resp"

curl -v "localhost:8081/body/1475276740074193887-4/req"
curl -v "localhost:8081/body/1475276740074193887-4/resp"
```


 
---------------------------------------------

[go]:http://golang.org/

