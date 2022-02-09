package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fergloragain/apitrace-remote/persistence"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/png"
	"net/http"
	"os"
	"strconv"
)

// Add a new Trace to the DB
func GetImage(traceDB *persistence.Cache) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		appName := p.ByName("name")
		imageID := p.ByName("image")

		val, err := traceDB.Get(fmt.Sprintf("%s", appName))

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`GetImage: could not find trace with appName <%s> and image ID <%s>
Error: %s`, appName, imageID, err.Error())))
			return
		}

		var trace Trace

		err = json.Unmarshal(val.([]byte), &trace)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetImage: could not unmarshal JSON for trace <%s>
Error: %s`, appName, err.Error())))
			return
		}

		imageFile, err := os.Open(fmt.Sprintf("%s/%s", trace.TargetDirectory, imageID))

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetImage: unable to open image <%s>
Error: %s`, imageFile, err.Error())))
			return
		}

		defer imageFile.Close()

		im, _, err := image.Decode(imageFile)

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetImage: unable to decode image <%s>
Error: %s`, imageFile, err.Error())))
			return
		}

		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, im); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetImage: unable to encode image as PNG <%s>
Error: %s`, imageFile, err.Error())))
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`GetImage: unable to write image <%s>
Error: %s`, imageFile, err.Error())))
			return
		}
	}
}
