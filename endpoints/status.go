package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Retrieve specific details about a particular app in the DB
func Status(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	w.Write([]byte("OK"))
	//w.WriteHeader(200)

}
