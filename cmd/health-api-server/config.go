package main

import "net/url"

type Configuration struct {
	BindAddress string
	DataPath    string
	ProxyURL    *url.URL
	HttpPath    string
}
