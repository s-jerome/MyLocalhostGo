package netflix

import (
	"bytes"
	"encoding/json"
	"io"
	"mylocalhost/logger"
	responses "mylocalhost/utils/responses"
	"net/http"
)

// I just added or removed a movie/serie from my playlist.
// My Chrome extension intercepted the request and sent the video data to be saved in database.
func SaveVideoToPlaylistRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The request must be POST")
		return
	}

	var requestBody, requestBodyErr = io.ReadAll(r.Body)
	if requestBodyErr != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, requestBodyErr, "Reading POST data")
		return
	}
	if len(requestBody) == 0 {
		responses.SendSimpleErrorMessageResponse(w, http.StatusBadRequest, "The POST data is empty")
		return
	}

	var videoData *videoData
	if parseErr := json.Unmarshal(requestBody, &videoData); parseErr != nil {
		var body = string(requestBody)
		logger.WriteError("[Netflix][SaveVideoToPlaylistRequestHandler] Body received:\n%s", body)
		responses.SendErrorResponse(w, http.StatusBadRequest, parseErr, "Parsing the POST data to JSON")
		return
	}

	var sqlResult = saveVideoToPlaylist(videoData)

	var buffer bytes.Buffer
	if encodeErr := json.NewEncoder(&buffer).Encode(sqlResult); encodeErr == nil {
		buffer.WriteTo(w)
	} else {
		responses.SendErrorResponse(w, http.StatusInternalServerError, encodeErr, "Encoding the SQL result in JSON")
	}
}
