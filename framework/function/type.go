package function

import "net/http"

// ScwFuncV1 is the prototype of a function handler that support std http objects.
// Version is embedded in the type to allow evolutions, this can allow core runtime to
// dynamically check for type to set approrite behaviour.
type ScwFuncV1 func(http.ResponseWriter, *http.Request)
