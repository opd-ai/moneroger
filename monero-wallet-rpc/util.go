package monerowalletrpc

import "net/url"

func validateRemoteDaemon(uri string) (scheme, host, port string, err error) {
	newUri, err := url.Parse(uri)
	if err != nil {
		return
	}
	host = newUri.Host
	port = newUri.Port()
	if port == "" {
		port = "18081"
	}
	return
}
