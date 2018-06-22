package ping

import (
	"fmt"
	"net/http"
)

// swagger:route GET /ping ping PingHandler
//
// Responses:
//   200:
//
// Produces:
// - text/plain
func PingHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, "pong\n")
}
