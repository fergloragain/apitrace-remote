package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
)

type App struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Executable  string   `json:"executable"`
	APITrace    string   `json:"apiTrace"`
	Retrace     string   `json:"retrace"`
	Timeout     int      `json:"timeout"`
	User        string   `json:"user"`
	PrivateKey  string   `json:"privateKey"`
	BuildScript string   `json:"buildScript"`
	Active      bool     `json:"active"`
	Branch      string   `json:"branch"`
	Traces      []string `json:"traces"`
	DumpImages  bool     `json:"dumpImages"`
}

type NewAppRequest struct {
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

type AppDescription struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Branch      string `json:"branch"`
	Traces      int    `json:"traces"`
}

// Get a list of descriptions for all apps in the appsDB
func GetApps(appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		// retrieve a list of all the app IDs in the appsDB
		allAppIDs := appsDB.TopLevelKeys()

		descriptionsArray := []AppDescription{}

		// for each appID, fetch the app data and compile a description
		for _, appID := range allAppIDs {
			appData, err := appsDB.Get(appID)

			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte(fmt.Sprintf(`GetApps: could not get application with ID <%s>
Error: %s`, appID, err.Error())))
				return
			}

			var app App

			err = json.Unmarshal(appData.([]byte), &app)

			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf(`GetApps: could not unmarshal application data with ID <%s>
Error: %s`, appID, err.Error())))
				return
			}

			descriptionsArray = append(descriptionsArray, AppDescription{
				app.ID,
				app.Name,
				app.Description,
				app.Branch,
				len(app.Traces),
			})
		}

		jsonResponse, err := json.Marshal(descriptionsArray)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetApps: could not marshal descriptions array
Error: %s`, err.Error())))
		} else {
			w.Write(jsonResponse)
		}
	}
}

// Add a new App to the DB
func AddApp(appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		name := p.ByName("name")

		var newAppRequest NewAppRequest

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddApp: could not read request body
Error: %s`, err.Error())))
			return
		}

		if err := r.Body.Close(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddApp: could not close body
Error: %s`, err.Error())))
			return
		}

		if err := json.Unmarshal(body, &newAppRequest); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddApp: could not unmarshal request body
Error: %s`, err.Error())))
			return
		}

		description := newAppRequest.Description
		url := newAppRequest.URL
		executable := newAppRequest.Executable
		apiTrace := newAppRequest.APITrace
		retrace := newAppRequest.Retrace
		user := newAppRequest.User
		privateKey := newAppRequest.PrivateKey
		buildScript := newAppRequest.BuildScript
		branch := newAppRequest.Branch
		timeout := newAppRequest.Timeout
		dumpImages := newAppRequest.DumpImages

		newID := appsDB.GetValidID(name)

		app := App{
			newID,
			name,
			description,
			url,
			executable,
			apiTrace,
			retrace,
			timeout,
			user,
			privateKey,
			buildScript,
			false,
			branch,
			[]string{},
			dumpImages,
		}

		applicationJSON, err := json.Marshal(app)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`AddApp: could not marshal application JSON
Error: %s`, err.Error())))
		} else {
			appsDB.Set(app.ID, applicationJSON)
			w.Write(applicationJSON)
		}
	}
}

// Retrieve specific details about a particular app in the DB
func GetApp(appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")

		applicationJSON, err := appsDB.Get(appName)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`GetApp: could not get application JSON for application name <%s>
Error: %s`, appName, err.Error())))
			return
		} else {
			w.Write(applicationJSON.([]byte))
		}
	}
}

// Update a particular app in the DB
func UpdateApp(appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")

		applicationJSON, err := appsDB.Get(appName)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not get application JSON for application name <%s>
Error: %s`, appName, err.Error())))
			return
		}

		var app App
		err = json.Unmarshal(applicationJSON.([]byte), &app)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not unmarshal application JSON for application name <%s>
Error: %s`, appName, err.Error())))
			return
		}

		var nar NewAppRequest

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not read request body
Error: %s`, appName, err.Error())))
			return
		}
		if err := r.Body.Close(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not close body reader
Error: %s`, appName, err.Error())))
			return
		}
		if err := json.Unmarshal(body, &nar); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not unmarshal JSON body
Error: %s`, appName, err.Error())))
			return
		}

		description := nar.Description
		url := nar.URL
		executable := nar.Executable
		apiTrace := nar.APITrace
		retrace := nar.Retrace
		user := nar.User
		privateKey := nar.PrivateKey
		buildScript := nar.BuildScript
		branch := nar.Branch
		timeout := nar.Timeout
		name := nar.Name
		dumpImages := nar.DumpImages

		updatedApplication := App{
			appName,
			name,
			description,
			url,
			executable,
			apiTrace,
			retrace,
			timeout,
			user,
			privateKey,
			buildScript,
			false,
			branch,
			app.Traces,
			dumpImages,
		}

		appJSON, err := json.Marshal(updatedApplication)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`UpdateApp: could not marshal updated application JSON
Error: %s`, appName, err.Error())))
		} else {
			appsDB.Set(appName, appJSON)
			w.Write(appJSON)
		}
	}

}

// Delete a particular app from the DB
func DeleteApp(appsDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(501)
		w.Write([]byte("DeleteApp not implemented"))
		//appName := p.ByName("name")
		//
		//appData, err := appsDB.Get(appName)
		//
		//if err != nil {
		//	log.Println(err.Error())
		//}
		//
		//var a App
		//err = json.Unmarshal(appData.([]byte), &a)
		//
		//if err != nil {
		//	log.Println(err.Error())
		//}
		//
		//traces := a.Traces
		//
		//for _, trace := range traces {
		//	tr, err := traceDB.Get(trace)
		//
		//	if err != nil {
		//		log.Println(err.Error())
		//	}
		//
		//	var t Trace
		//	err = json.Unmarshal(tr.([]byte), &t)
		//
		//	if err != nil {
		//		log.Println(err.Error())
		//	}
		//
		//	rets := t.Retraces
		//
		//	for _, retrace := range rets {
		//
		//		retraceDB.Delete(retrace)
		//
		//	}
		//
		//	for i := 0; i < t.NumberOfFrames; i++ {
		//
		//		frameID := fmt.Sprintf("%s-%d", t.ID, i)
		//
		//		dumpDB.Delete(frameID)
		//	}
		//
		//	traceDB.Delete(trace)
		//
		//	os.RemoveAll(t.TargetDirectory)
		//
		//}
		//
		//err = appsDB.Delete(appName)
		//
		//if err != nil {
		//	w.Write([]byte(err.Error()))
		//} else {
		//	w.Write([]byte("OK"))
		//}
	}

}
