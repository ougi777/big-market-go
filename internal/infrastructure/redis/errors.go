package redis

import "errors"

var ErrClientNotConnected = errors.New("redis client is not connected")
