package main

import (
	"encoding/gob"
	"github.com/dgraph-io/badger"
	"github.com/fergloragain/apitrace-remote/endpoints"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	opts := badger.DefaultOptions
	opts.Dir = "./db"
	opts.ValueDir = "./db"
	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	gob.Register(map[string][]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register([]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register([][]int{})

	appsDB := persistence.NewCache(db, "apps")
	traceDB := persistence.NewCache(db, "trace")
	dumpDB := persistence.NewCache(db, "dump")
	retraceDB := persistence.NewCache(db, "retrace")
	configDB := persistence.NewCache(db, "config")

	router := httprouter.New()

	router.GET("/apps", endpoints.GetApps(appsDB))
	router.POST("/apps/:name", endpoints.AddApp(appsDB))
	router.GET("/apps/:name", endpoints.GetApp(appsDB))
	router.PUT("/apps/:name", endpoints.UpdateApp(appsDB))
	router.DELETE("/apps/:name", endpoints.DeleteApp(appsDB))

	router.GET("/traces", endpoints.GetTraces(traceDB))
	router.POST("/traces/:name", endpoints.AddTrace(traceDB, appsDB, dumpDB))
	router.GET("/traces/:name", endpoints.GetTrace(traceDB))
	router.DELETE("/traces/:name", endpoints.DeleteTrace(traceDB))

	router.GET("/dumps/:name/:frame", endpoints.GetDump(dumpDB))
	router.DELETE("/dumps/:name", endpoints.DeleteDump(dumpDB, traceDB))

	router.POST("/retrace/:name/:call", endpoints.AddRetrace(retraceDB, traceDB, appsDB))
	router.GET("/retrace/:name/:call", endpoints.GetRetrace(retraceDB))

	router.GET("/images/:name/:image", endpoints.GetImage(traceDB))

	router.GET("/config", endpoints.GetConfig(configDB))
	router.PUT("/config", endpoints.UpdateConfig(configDB))

	router.POST("/killmenow", endpoints.Kill)
	router.GET("/status", endpoints.Status)

	err = http.ListenAndServe(":8080", &Server{router})

	if err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	r *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	s.r.ServeHTTP(w, r)
}
