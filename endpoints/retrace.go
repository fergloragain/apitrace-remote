package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/fergloragain/apitrace-remote/operations"
	"github.com/fergloragain/apitrace-remote/parsers"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Retrace struct {
	ID              string              `json:"id"`
	AppID           string              `json:"appID"`
	Name            string              `json:"name"`
	Status          string              `json:"status"`
	ImageDumpStdout string              `json:"imageDumpStdout"`
	ImageDumpStderr string              `json:"imageDumpStderr"`
	RetraceStdout   string              `json:"retraceStdout"`
	RetraceStderr   string              `json:"retraceStderr"`
	RetraceData     parsers.RetraceData `json:"retraceData"`
	ImageSet        *parsers.ImageSet   `json:"imageSet"`
}

// Add a new Trace to the DB
func AddRetrace(retraceDB, traceDB, appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")
		callID := p.ByName("call")

		val, err := traceDB.Get(appName)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Unable to retrieve information for %s: %s", appName, err.Error())))
			return
		}

		var trace Trace
		err = json.Unmarshal(val.([]byte), &trace)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Unable to unmarshal data for %s: %s", appName, err.Error())))
			return
		}

		val, err = appsDB.Get(trace.AppID)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Unable to retrieve information for %s: %s", appName, err.Error())))
			return
		}

		var app App
		err = json.Unmarshal(val.([]byte), &app)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Unable to unmarshal data for %s: %s", appName, err.Error())))
			return
		}

		retraceID := fmt.Sprintf("%s-%s", appName, callID)

		val, err = retraceDB.Get(retraceID)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Unable to retrieve retrace information for %s: %s", appName, err.Error())))
			return
		}

		if val != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("A retrace for %s-%s already exists", appName, callID)))
			return
		}

		retraceStatus := Retrace{
			retraceID,
			trace.AppID,
			retraceID,
			Pending,
			"",
			"",
			"",
			"",
			parsers.RetraceData{},
			nil,
		}

		retraceStatusJSON, err := json.Marshal(retraceStatus)

		if err != nil {
			w.Write([]byte(fmt.Sprintf("Unable to marshal trace data for %s: %s", appName, err.Error())))
			return
		}

		retraceDB.Set(retraceID, retraceStatusJSON)

		trace.Retraces = append(trace.Retraces, retraceID)

		traceJSON, err := json.Marshal(trace)

		if err != nil {
			w.Write([]byte(fmt.Sprintf("Unable to marshal trace data for %s: %s", appName, err.Error())))
			return
		}

		// update the app in the appsDB to register it as as having an active job
		traceDB.Set(trace.ID, traceJSON)

		// now, kick off the job asynchronously
		go func() {

			// trace the application
			retraceStdout, retraceStderr, err := operations.Retrace(trace.TargetDirectory, app.Retrace, trace.TraceFile, callID)

			if app.DumpImages {
				log.Println("Dumping images...")

				imageDumpStdout, imageDumpStderr, err := operations.DumpImages(trace.TargetDirectory, app.APITrace, trace.TraceFile, callID)

				log.Println("Finished dumping images...")

				if err != nil {
					log.Println(err.Error())
				}

				traceImageDump := parsers.ParseImageDumpFile(imageDumpStdout)
				retraceStatus.ImageSet = traceImageDump
				retraceStatus.ImageDumpStdout = imageDumpStdout
				retraceStatus.ImageDumpStderr = imageDumpStderr
			}

			log.Println("FINISHED TRACE")
			// once all processes have finished, mark the trace status as complete, and save in the build and trace outputs
			retraceStatus.Status = Complete
			retraceStatus.RetraceStdout = retraceStdout
			retraceStatus.RetraceStderr = retraceStderr

			retraceStructs := parsers.NewRetraceData(retraceStdout)

			retraceStatus.RetraceData = retraceStructs

			retraceJSON, err := json.Marshal(retraceStatus)

			if err != nil {
				log.Println("Error marshalling retraceStatus")
				log.Println(err.Error())
			}

			retraceDB.Set(retraceID, retraceJSON)

		}()

		w.Write(retraceStatusJSON)
	}

}

// Add a new Trace to the DB
func GetRetrace(retraceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")
		callID := p.ByName("call")

		val, err := retraceDB.Get(fmt.Sprintf("%s-%s", appName, callID))

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`GetRetrace: Unable to retrieve information for <%s>
Error: %s`, appName, err.Error())))
			return
		}

		w.Write(val.([]byte))
	}

}
