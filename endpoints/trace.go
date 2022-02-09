package endpoints

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fergloragain/apitrace-remote/operations"
	"github.com/fergloragain/apitrace-remote/parsers"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	Pending  = "Pending"
	Complete = "Complete"
)

type Trace struct {
	ID              string   `json:"id"`
	AppID           string   `json:"appID"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	BuildStdout     string   `json:"buildStdout"`
	BuildStderr     string   `json:"buildStderr"`
	TraceStdout     string   `json:"traceStdout"`
	TraceStderr     string   `json:"traceStderr"`
	CloneStdout     string   `json:"cloneStdout"`
	CloneStderr     string   `json:"cloneStderr"`
	DumpStderr      string   `json:"dumpStderr"`
	TargetDirectory string   `json:"targetDirectory"`
	NumberOfFrames  int      `json:"numberOfFrames"`
	Retraces        []string `json:"retraces"`
	TraceFile       string   `json:"traceFile"`
}

// Get a list of all the app IDs within the DB
func GetTraces(traceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		allTraces := traceDB.TopLevelKeys()

		tracesJSON, err := json.Marshal(allTraces)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetTraces: could not marshal all traces
Error: %s`, err.Error())))
			return
		}

		w.Write(tracesJSON)
	}

}

// Add a new Trace to the DB
func AddTrace(traceDB, appsDB, dumpDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")

		val, err := appsDB.Get(appName)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddTrace: Unable to retrieve information for <%s>
Error: %s`, appName, err.Error())))
			return
		}

		var app App
		err = json.Unmarshal(val.([]byte), &app)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddTrace: Unable to unmarshal app data
Error: %s`, appName, err.Error())))
			return
		}

		if app.Active {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf("Build for %s is already actvie, please wait", appName)))
			return
		}

		potentialTraceID := fmt.Sprintf("%s-trace", app.ID)
		traceID := traceDB.GetValidID(potentialTraceID)

		traceStatus := Trace{
			ID:              traceID,
			AppID:           app.ID,
			Name:            traceID,
			Status:          Pending,
			BuildStdout:     "",
			BuildStderr:     "",
			TraceStdout:     "",
			TraceStderr:     "",
			TargetDirectory: "",
			NumberOfFrames:  0,
			Retraces:        []string{},
		}

		traceStatusJSON, err := json.Marshal(traceStatus)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddTrace: Unable to marshal trace status data
Error: %s`, appName, err.Error())))
			return
		}

		traceDB.Set(traceID, traceStatusJSON)

		// mark the app as having an active job, to prevent multiple jobs running at once
		app.Active = true

		app.Traces = append(app.Traces, traceID)

		appJSON, err := json.Marshal(app)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddTrace: Unable to marshal app data
Error: %s`, appName, err.Error())))
			return
		}

		// update the app in the appsDB to register it as as having an active job
		appsDB.Set(app.ID, appJSON)

		// now, kick off the job asynchronously
		go func() {

			// create a random target directory on the server
			targetDirectory := fmt.Sprintf("/tmp/%s-%d", app.ID, time.Now().Nanosecond())

			// clone the repo to the random directory
			cloneStdout, cloneStderr, err := operations.Clone(app.User, app.PrivateKey, app.URL, targetDirectory, app.Branch)

			fmt.Println("Clone finished")

			if err != nil {
				log.Println(fmt.Sprintf("AddTrace: Error cloning repo %s: %s", app.URL, err.Error()))
				return
			}

			// build the application
			buildStdout, buildStderr, err := operations.Build(targetDirectory, app.BuildScript)

			fmt.Println("Build finished")

			if err != nil {
				log.Println(fmt.Sprintf("AddTrace: Error building application %s: %s", app.Name, err.Error()))
				return
			}


			fmt.Println("buildStdout")
			fmt.Println(buildStdout)


			fmt.Println("buildStderr")
			fmt.Println(buildStderr)

			// trace the application
			traceStdout, traceStderr, err := operations.Trace(targetDirectory, app.APITrace, app.Executable, app.Timeout)

			fmt.Println("Trace finished")

			if err != nil {
				log.Println(fmt.Sprintf("AddTrace: Error tracing application %s: %s", app.Name, err.Error()))
				return
			}

			fmt.Println("traceStdout")
			fmt.Println(traceStdout)


			fmt.Println("traceStderr")
			fmt.Println(traceStderr)


			// get the tracefile name
			traceFile := getTraceFile(traceStderr)

			fmt.Println("Tracefile is")
			fmt.Println(traceFile)


			// dump the trace file
			dumpStdout, dumpStderr, err := operations.Dump(targetDirectory, app.APITrace, traceFile)

			fmt.Println("dumpStdout")
			fmt.Println(dumpStdout)


			fmt.Println("dumpStderr")
			fmt.Println(dumpStderr)


			if err != nil {
				log.Println(fmt.Sprintf("AddTrace: Error dumping trace %s: %s", traceFile, err.Error()))
				return
			}

			// once all processes have finished, mark the trace status as complete, and save in the build and trace outputs
			traceStatus.Status = Complete
			traceStatus.BuildStdout = buildStdout
			traceStatus.BuildStderr = buildStderr
			traceStatus.TraceStdout = traceStdout
			traceStatus.TraceStderr = traceStderr
			traceStatus.CloneStdout = cloneStdout
			traceStatus.CloneStderr = cloneStderr
			traceStatus.DumpStderr = dumpStderr
			traceStatus.TraceFile = traceFile

			traceStatus.TargetDirectory = targetDirectory

			traceDump := parsers.ParseDump(dumpStdout)

			fmt.Println("ParseDump finished")
			//fmt.Println(traceDump)


			// since timout kills the trace, the last frame will probably always be only partially complete, so we want to drop it from the frame collection
			if len(traceDump.Frames) > 0 {
				traceDump.Frames = traceDump.Frames[:len(traceDump.Frames)-1]
			}

			traceStatus.NumberOfFrames = len(traceDump.Frames)

			for i, frame := range traceDump.Frames {

				frameID := fmt.Sprintf("%s-%d", traceID, i)

				dumpFrame, err := json.Marshal(frame)

				if err != nil {
					log.Println(fmt.Sprintf("AddTrace: Error marshalling frame %d: %s", i, err.Error()))
					return
				}

				dumpDB.Set(frameID, dumpFrame)
			}

			updatedTraceStatusJSON, err := json.Marshal(traceStatus)

			if err != nil {
				w.Write([]byte(fmt.Sprintf("AddTrace: Unable to marshal trace status %+v", traceStatus)))
				return
			}

			traceDB.Set(traceID, updatedTraceStatusJSON)

			// mark the app as no longer having an active job
			app.Active = false

			updatedAppJSON, err := json.Marshal(app)

			if err != nil {
				w.Write([]byte(fmt.Sprintf("AddTrace: Unable to marshal app data for %+v", app)))
				return
			}

			appsDB.Set(app.ID, updatedAppJSON)

		}()

		w.Write(traceStatusJSON)
	}

}

// Retrieve specific details about a particular trace in the DB
func GetTrace(traceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		traceName := p.ByName("name")

		val, err := traceDB.Get(traceName)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`GetTrace: could not find trace with ID: <%s>
Error: %s`, traceName, err.Error())))
			return
		} else {
			w.Write(val.([]byte))
		}
	}

}

// Delete a particular trace from the DB
func DeleteTrace(traceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(501)
		w.Write([]byte("DeleteTrace not implemented"))
	}
	//traceName := p.ByName("name")
	//
	//traceData, err := traceDB.Get(traceName)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//var trace Trace
	//err = json.Unmarshal(traceData.([]byte), &trace)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//appData, err := appsDB.Get(trace.AppID)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//var app App
	//err = json.Unmarshal(appData.([]byte), &app)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//app.Traces = remove(app.Traces, traceName)
	//
	//updatedAppJSON, err := json.Marshal(app)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//
	//appsDB.Set(app.ID, updatedAppJSON)
	//
	//for i := 0; i < trace.NumberOfFrames; i++ {
	//	dumpID := fmt.Sprintf("%s-%d", traceName, i)
	//
	//	err = dumpDB.Delete(dumpID)
	//
	//	if err != nil {
	//		w.Write([]byte(err.Error()))
	//	}
	//}
	//
	//err = traceDB.Delete(traceName)
	//
	//if err != nil {
	//	w.Write([]byte(err.Error()))
	//} else {
	//	w.Write([]byte("OK"))
	//}

}

func remove(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}

	return slice
}

func getTraceFile(traceStdErr string) string {

	traceFilePath := ""
	scanner := bufio.NewScanner(strings.NewReader(traceStdErr))

	for scanner.Scan() {
		line := scanner.Text()

		stringFields := strings.Fields(line)

		if strings.HasPrefix(line, "apitrace:") {
			if strings.Contains(line, "tracing to") {
				traceFilePath = stringFields[len(stringFields)-1]
				break
			}
		}

	}

	return traceFilePath

}
