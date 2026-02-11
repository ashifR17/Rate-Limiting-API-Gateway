package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Safe Director function
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		// Strip "/api" prefix if present
		if strings.HasPrefix(req.URL.Path, "/api") {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
