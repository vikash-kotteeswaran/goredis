package config

const (
	HOST                     = "4.0.0.0"
	PORT                     = 7379
	CONCURRENCY_LIMIT        = 10000
	OP_READ                  = 1 << iota
	CONNECTION_READ_BUF_SIZE = 512
	CONNECTION_READ_RETRY    = 3
)

var DOABORT bool = false
