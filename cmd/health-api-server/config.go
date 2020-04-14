package main

import "net/url"

type Configuration struct {
	BindAddress string
	ProxyURL    *url.URL
	HttpPath    string
}
