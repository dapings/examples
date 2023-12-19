package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type Func func(r *http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r roundRobinSwitcher) GetProxy(_ *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]

	return u, nil
}

// RoundRobinProxySwitcher creates a proxy switcher function
// which rotates proxyURLs on every request.
// The proxy type is determined by the URL scheme.
// "http", "https" and "socks5" are supported.
// If the scheme is empty, "http" is assumed.
func RoundRobinProxySwitcher(proxyURLs ...string) (Func, error) {
	if len(proxyURLs) < 1 {
		return nil, errors.New("proxy URL list is empty")
	}

	urls := make([]*url.URL, len(proxyURLs))

	for i, u := range proxyURLs {
		parsedURL, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		urls[i] = parsedURL
	}

	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}
