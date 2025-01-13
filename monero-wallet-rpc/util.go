package monerowalletrpc

import (
	"fmt"
	"net/url"
)

// validateRemoteDaemon ensures that the remote daemon URL is correctly formed
// and outputs the components of the URL.
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
	if len(newUri.Path) > 1 {
		err = fmt.Errorf("Remote node URLs may not contain a path: %s", newUri.Path)
		return
	}
	return
}
