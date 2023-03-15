package builder

import (
	"net/url"
	"time"
)

type Node struct {
	CACert            string
	Cert              string
	Key               string
	Address           *url.URL
	ConnectionTimeout time.Duration
}
