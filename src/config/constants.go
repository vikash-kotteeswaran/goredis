package config

const (
	HOST              = "0.0.0.0"
	PORT              = 7379
	CONCURRENCY_LIMIT = 10000
	OP_READ           = 1 << iota
)

var DOABORT bool = false
