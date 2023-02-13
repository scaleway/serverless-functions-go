package server

import (
	"fmt"
)

// ErrPayloadTooLarge - Error type for payload size is grater that anticipated.
var ErrPayloadTooLarge = fmt.Errorf("request payload too large, max payload size = %d bytes", payloadMaxSize)
