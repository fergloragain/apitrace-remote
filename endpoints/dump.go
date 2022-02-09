package endpoints

import (
	"fmt"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Retrieve specific details about a particular dump in the DB
func GetDump(dumpDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		dumpName := p.ByName("name")
		frameNumber := p.ByName("frame")

		fn, err := strconv.Atoi(frameNumber)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetDump: could not convert <%s> to an int
Error: %s`, frameNumber, err.Error())))
			return
		}

		dumpID := fmt.Sprintf("%s-%d", dumpName, fn)

		val, err := dumpDB.Get(dumpID)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`GetDump: could not get dump for ID <%s>
Error: %s`, dumpID, err.Error())))
		} else {
			w.Write(val.([]byte))
		}
	}

}

// Delete a particular dump from the DB
func DeleteDump(dumpDB, traceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(501)
		w.Write([]byte("DeleteDump not implemented"))
	}
	//dumpName := p.ByName("name")
	//
	//trace, err := traceDB.Get(dumpName)
	//
	//if err != nil {
	//	w.Write([]byte(fmt.Sprintf("Error fetching trace for %s", dumpName)))
	//	return
	//}
	//
	//var t Trace
	//err = json.Unmarshal(trace.([]byte), &t)
	//
	//if err != nil {
	//	w.Write([]byte(fmt.Sprintf("Error unmarshalling trace for %s", dumpName)))
	//	return
	//}
	//
	//for i := 0; i < t.NumberOfFrames; i++ {
	//
	//	dumpID := fmt.Sprintf("%s-%d", dumpName, i)
	//	err = dumpDB.Delete(dumpID)
	//
	//	if err != nil {
	//		w.Write([]byte(err.Error()))
	//	}
	//}
	//
	//w.Write([]byte("OK"))

}
