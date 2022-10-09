package data

import (
	"net/http"
	"time"
)

// Client that used to do all requests from this package
var httpClient = http.Client{
	Timeout: 10 * time.Second,
}
