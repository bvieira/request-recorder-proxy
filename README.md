 ---------------------------------------------
## Request Recorder Proxy
tool using [go] to help automated tests validations for asynchronous services

### Flags:
Use '-h' for flags usage
```sh
request-recorder-proxy -h
```

### Services:
	- ?

### Example:

Request using proxy
```sh
curl "http://www.google.com/" --proxy localhost:8080 -H "X-tid: test"
```

Recover infos
```sh
curl -v "localhost:8081/metadata/req?key=test&uri=www.google.com/&method=get"
curl -v "localhost:8081/metadata/resp?key=test&uri=www.google.com/&method=get"
```


 
---------------------------------------------

[go]:http://golang.org/

