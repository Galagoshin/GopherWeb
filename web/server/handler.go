package server

import (
	"fmt"
	"github.com/Galagoshin/GoLogger/logger"
	"github.com/Galagoshin/GoUtils/requests"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Error(err)
	}
	found := false
	for _, route := range routes {
		if route.pattern.MatchString(r.RequestURI) {
			allowed := false
			for _, allowedMethod := range route.methods {
				if requests.Method(r.Method) == allowedMethod {
					allowed = !allowed
					break
				}
			}
			if allowed {
				data := r.Form
				if route.option.len() > 0 {
					point := false
					name, val, skip, uris := "", "", "", 0
					for _, char := range []rune(route.option) {
						if string(char) == "{" && !point {
							point = true
							in := strings.Index(string(route.option)[len(skip):], "}")
							sep := "/"
							if in+1 < len(string(route.option)[len(skip):]) {
								sep = string(route.option)[len(skip):][in+1 : in+2]
							}
							val = strings.Split(r.RequestURI[uris:], sep)[0]
							uris += len(val)
							skip += string(char)
						} else if string(char) == "}" && point {
							data.Add(name, val)
							skip += name + "}"
							name = ""
							val = ""
							point = false
						} else if point {
							name += string(char)
						} else {
							skip += string(char)
							uris++
						}
					}
				}
				request := &requests.Request{
					Method:  requests.Method(r.Method),
					Data:    data,
					Url:     requests.URL(r.URL.String()),
					Headers: r.Header,
					Cookies: r.Cookies(),
				}
				response, err := route.handle(request)
				if err != nil {
					logger.Print(fmt.Sprintf("%s: %s (%d)", r.Method, r.RequestURI, 500))
					//call event cancel
					w.WriteHeader(500)
					_, err := w.Write([]byte("500 internal server error"))
					if err != nil {
						logger.Error(err)
					}
					return
				}
				status_code := 200
				if response.StatusCode != 0 {
					status_code = response.StatusCode
				}
				for _, cookie := range response.Cookies {
					http.SetCookie(w, cookie)
				}
				//call event cancel
				for key, val := range response.Header {
					w.Header().Add(key, val[0])
				}
				w.WriteHeader(status_code)
				_, err = w.Write(response.Body)
				if err != nil {
					logger.Error(err)
				}
				logger.Print(fmt.Sprintf("%s: %s (%d)", r.Method, r.RequestURI, status_code))
				logger.Debug(5, false, fmt.Sprintf("Request data: %+v", request.Data))
				logger.Debug(6, false, fmt.Sprintf("Request headers: %+v", request.Headers))
				logger.Debug(6, false, fmt.Sprintf("Request cookies: %+v", request.Cookies))
				logger.Debug(6, false, fmt.Sprintf("Response body: %+v", string(response.Body)))
				logger.Debug(6, false, fmt.Sprintf("Response headers: %+v", response.Header))
				logger.Debug(6, false, fmt.Sprintf("Response cookies: %+v", response.Cookies))
			} else {
				logger.Print(fmt.Sprintf("%s: %s (%d)", r.Method, r.RequestURI, 405))
				//call event cancel
				w.WriteHeader(405)
				_, err := w.Write([]byte("405 method not allowed"))
				if err != nil {
					logger.Error(err)
				}
			}
			found = true
			break
		}
	}
	if !found {
		logger.Print(fmt.Sprintf("%s: %s (%d)", r.Method, r.RequestURI, 404))
		//call event cancel
		w.WriteHeader(404)
		_, err := w.Write([]byte("404 not found"))
		if err != nil {
			logger.Error(err)
		}
	}
}
