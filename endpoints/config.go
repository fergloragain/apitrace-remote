package endpoints

import (
	"encoding/json"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Config struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Executable  string `json:"executable"`
	APITrace    string `json:"apiTrace"`
	Retrace     string `json:"retrace"`
	User        string `json:"user"`
	PrivateKey  string `json:"privateKey"`
	BuildScript string `json:"buildScript"`
	Branch      string `json:"branch"`
	Timeout     int    `json:"timeout"`
	DumpImages  bool   `json:"dumpImages"`
}

// Add a new App to the DB
func SetConfig(configDB *persistence.Cache, body []byte) {

	var c Config

	err := json.Unmarshal(body, &c)

	if err != nil {
		log.Print(err.Error())
	}

	configJSON, err := json.Marshal(c)

	if err != nil {
		log.Print(err.Error())
	}

	configDB.Set("config", configJSON)
}

// Retrieve specific details about a particular app in the DB
func GetConfig(configDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(501)
		w.Write([]byte("GetConfig not implemented"))
	}
	//val, err := configDB.Get("config")
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//} else {
	//	w.Write(val.([]byte))
	//}
}

// Update a particular app in the DB
func UpdateConfig(configDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(501)
		w.Write([]byte("UpdateConfig not implemented"))
	}

	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//if err := r.Body.Close(); err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}

	//SetConfig(body)
}
