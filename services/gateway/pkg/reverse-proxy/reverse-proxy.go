package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func New(target string) http.Handler {
	url, _ := url.Parse(target)

	return httputil.NewSingleHostReverseProxy(url)
}
