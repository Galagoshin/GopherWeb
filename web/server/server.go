package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/configs"
	"net/http"
	time2 "time"
)

var Config = &configs.Config{Name: "server"}

var (
	host          string = "0.0.0.0"
	port          uint16 = 80
	domain        string
	checkHostname bool = false
	sslEnabled    bool = false
	crtFilename   string
	keyFilename   string
)

var srv *http.Server

var isHandlerInit = false

func Init() {
	Config.Init(map[string]any{
		"host":           "0.0.0.0",
		"port":           80,
		"domain":         "localhost",
		"check-hostname": "false",
		"enable-ssl":     "false",
		"crt-file":       "file.crt",
		"key-file":       "file.key",
	}, 1)

	host_init, hostError := Config.Get("host")
	if !hostError {
		logger.Panic(errors.New("\"host\" is not defined in the server config."))
	} else {
		host = host_init.(string)
	}
	port_init, portError := Config.Get("port")
	if !portError {
		logger.Panic(errors.New("\"port\" is not defined in the server config."))
	} else {
		port = uint16(port_init.(int))
	}
	domain_init, domainError := Config.Get("domain")
	if !domainError {
		logger.Panic(errors.New("\"domain\" is not defined in the server config."))
	} else {
		domain = domain_init.(string)
	}
	checkHostname_init, checkHostnameError := Config.Get("check-hostname")
	if !checkHostnameError {
		logger.Panic(errors.New("\"check-hostname\" is not defined in the server config."))
	} else {
		checkHostname = checkHostname_init.(string) == "true"
	}
	enableSsl_init, enableSslError := Config.Get("enable-ssl")
	if !enableSslError {
		logger.Panic(errors.New("\"enable-ssl\" is not defined in the server config."))
	} else {
		sslEnabled = enableSsl_init.(string) == "true"
	}
	if enableSsl_init.(string) == "true" {
		crtFilenameKey, crtFilenameError := Config.Get("crt-file")
		if !crtFilenameError {
			logger.Panic(errors.New("\"crt-file\" is not defined in the server config."))
		} else {
			crtFilename = crtFilenameKey.(string)
		}
		keyFilenameKey, keyFilenameError := Config.Get("key-file")
		if !keyFilenameError {
			logger.Panic(errors.New("\"key-file\" is not defined in the server config."))
		} else {
			keyFilename = keyFilenameKey.(string)
		}
	}
	srv = &http.Server{Addr: fmt.Sprintf("%s:%d", host, port)}

	if !isHandlerInit {
		isHandlerInit = !isHandlerInit
		http.HandleFunc("/", handler)
	}
}

func Stop() {
	logger.Print("Stopping webserver...")
	if err := srv.Shutdown(context.TODO()); err != nil {
		logger.Panic(err)
	}
	time2.Sleep(time2.Second)
	logger.Print("Webserver stopped.")
}

func Run() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r.(error))
			Run()
		}
	}()

	logger.Print(fmt.Sprintf("Webserver starting on %s:%d", host, port))
	if sslEnabled {
		err := srv.ListenAndServeTLS(crtFilename, keyFilename)
		if err != http.ErrServerClosed {
			logger.Panic(err)
		}
	} else {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Panic(err)
		}
	}
}
