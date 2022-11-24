# GopherWeb
Open-source, high-performance web framework for the Golang.

### Contents:
- [Getting started](#getting-started)
    - [Install](#install)
    - [Example usage](#example-usage)
    - [Web server configuration](#web-server-configuration)
    - [Framework configuration](#framework-configuration)
    - [Developer mode](#developer-mode)
    - [Live reloads](#live-reloads)
    - [Hot reload](#hot-reload)
- [Routing](#routing)
    - [Example route](#example-route)
    - [Parametric routes](#parametric-routes)
    - [Static files](#static-files)
- [Render core](#render-core)
    - [Standart partials](#standart-partials)
    - [HTML escaped partials](#html-escaped-partials)
    - [Includes](#includes)

## Getting started

### Install

`go get -t github.com/Galagoshin/GopherWeb`

### Example usage

```go
package main

import (
	"fmt"
	"github.com/Galagoshin/GoUtils/requests"
	"github.com/Galagoshin/GopherWeb/web"
	"github.com/Galagoshin/GopherWeb/web/server"
)

func main(){
	web.Init()
	server.Route("/", []requests.Method{requests.GET}, index)
	web.Run()
}

func index(request *requests.Request) (requests.Response, error) {
	val := request.Data.Get("field")
	return requests.Response{
		Body: []byte(fmt.Sprintf("field value is %s", val)),
	}, nil
}
```

After the first launch, GopherWeb will create configuration files.

### Web server configuration

#### server.gconf
```yaml
#v=1
host: 0.0.0.0
port: 80
domain: localhost
check-hostname: false
enable-ssl: false
crt-file: file.crt
key-file: file.key
```

* `host` - The IP address where the web server will run
* `port` - The port where the web server will run
* `domain` - Server domain name
* `check-hostname` - Allows accessing the webserver only with the hostname from the `domain` parameter
* `enable-ssl` - Enable SSL
* `crt-file` - Path to crt file (works with `enable-ssl: false`)
* `key-file` - Path to key file (works with `enable-ssl: false`)

### Framework configuration

#### gopher.gconf

```yaml
#v=1
live-reload-enabled: true
hot-reload-enabled: true
write-logs: true
debug-level: 0
```

* `live-reload-enabled` - Enable live reloads
* `hot-reload-enabled` - Enable hot reload
* `write-logs` - Write logs
* `debug-level` - Set debug level

### Developer mode

Developer mode is disabled automatically when you run your application inside a docker container.

To force disable developer mode use the flag `--mode prod`

You can also enable developer mode inside a docker container with the flag `--mode dev`

### Live reloads

GopherWeb supports live reloads in developer mode. You need to correctly configure the GoLand configuration and `build.gconf`.

#### build.gconf
```yaml
#v=1
build: build -o gopher_server src/main.go
run: run src/main.go
```

The `run` parameter should run your application using standart `go` utility.

#### GoLand configuration

![Example GoLand configuration](https://github.com/Galagoshin/GopherWeb/blob/master/example_conf.jpg?raw=true)

You need to run `main()` in `github.com/Galagoshin/GopherWeb/develop/main.go`

### Hot reload

GopherWeb supports hot reload. If you change any configuration file or content in the `views` directory, the web server will be restarted with the new settings.

## Routing

### Example route

```go
server.Route("/api/info", []requests.Method{requests.GET, requests.POST}, info)

func info(request *requests.Request) (requests.Response, error) {
  if request.Method == requests.POST {
    //Create info
    return requests.Response{
      Body: []byte("created info."),
      StatusCode: http.StatusCreated,
    }, nil
  }else if request.Method == requests.GET{
    return requests.Response{
      Body: []byte("any info."),
      StatusCode: 200,
    }, nil
  }else{
    return requests.Response{
      Body: []byte("method not allowed."),
      StatusCode: http.StatusMethodNotAllowed,
    }, nil
  }
}
```

### Parametric routes

#### Example usage

```go
server.Route("/api/user/{id: [0-9]}", []requests.Method{requests.GET}, user)

func user(request *requests.Request) (requests.Response, error) {
  user_id := request.Data.Get("id")
  return requests.Response{
    Body: []byte(fmt.Sprintf("getting user id %s", user_id)),
  }, nil
}
```

#### Example without parameter

```go
server.Route("/api/{[A-Za-z0-9]}/get", []requests.Method{requests.GET}, hander)
```

### Static files

All static files are located in the `static/` directory and are available via the `/static` route.

## Render core

### Standart partials

`index.html` file in `views/` directory.
```html
<head>
    <title><%=title%></title>
</head>
```

Usage in application.
```go
func any_handler(request *requests.Request) (requests.Response, error) {
  response := requests.Response{
    Body: []byte(render.GetView("index").Render(map[string]string{
      "title": "Page title",
    })),
    StatusCode: 200,
  }
  return response, nil
}
```

### HTML escaped partials

`index.html` file in `views/` directory.
```html
<head>
    <title><%!title%></title>
</head>
```

Usage in application.
```go
func any_handler(request *requests.Request) (requests.Response, error) {
  response := requests.Response{
    Body: []byte(render.GetView("index").Render(map[string]string{
      "title": "<script>alert(1);</script>",
    })),
    StatusCode: 200,
  }
  return response, nil
}
```

### Includes

In this example, the content of `views/path/to/file.html` is set to the current html file.
```html
<head>
    <title><%& path/to/file %></title>
</head>
```

