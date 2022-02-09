package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

func Kill(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	os.Exit(1987)
}
