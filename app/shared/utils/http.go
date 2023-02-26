package utils

import (
	"errors"
	"net/http"
	"net/url"
)

const (
	uriVerifyVPN = "https://api.notion.com/v1"
	errVPN       = "Get \"https://api.notion.com/v1\": dial tcp: lookup api.notion.com: no such host"
)

// GetClient return a client depending on if needed the proxy or not
func GetClient(proxy string) *http.Client {
	httpClient := &http.Client{}
	_, err := httpClient.Get(uriVerifyVPN)

	if err == nil {
		return httpClient
	}

	if err.Error() == errVPN {
		transport, errGetHTTPTransport := getHTTPTransport(proxy)
		if errGetHTTPTransport != nil {
			panic(errGetHTTPTransport)
		}

		httpClient.Transport = transport
		/**/
		_, err2 := httpClient.Get(uriVerifyVPN)

		if err2 == nil {
			return httpClient
		}

		return httpClient
	}

	panic(err)
}

func getHTTPTransport(urlHTTPProxy string) (*http.Transport, error) {
	if urlHTTPProxy == "" {
		return nil, errors.New("var HTTP_PROXY not set")
	}

	proxyUrl, err := url.Parse(urlHTTPProxy)
	if err != nil {
		return nil, err
	}
	return &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, nil
}
