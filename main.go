package shorturl

import (
	"net/http"

	"salbo.ai/short-url/handler"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    handler.HandleRequest(w, r)
}
