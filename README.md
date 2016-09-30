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

### Version
/version

### Requests
/requests?key=&uri=&method=

### Metadata
/metadata/:id/:type

### Body
/body/:id/:type


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

