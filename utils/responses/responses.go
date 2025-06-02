package utils

import (
	"bytes"
	"encoding/json"
	"mylocalhost/logger"
	"net/http"
)

func SendErrorResponse(w http.ResponseWriter, statusCode int, err error, operation string) {
	if err != nil {
		SendErrorMessageResponse(w, statusCode, err.Error(), operation)
	} else {
		SendErrorMessageResponse(w, statusCode, "", operation)
	}
}

func SendSimpleErrorMessageResponse(w http.ResponseWriter, statusCode int, errMessage string) {
	SendErrorMessageResponse(w, statusCode, errMessage, "")
}

func SendErrorMessageResponse(w http.ResponseWriter, statusCode int, errMessage string, operation string) {
	var data = make(map[string]string)
	if errMessage != "" {
		data["error"] = errMessage
	}
	if operation != "" {
		data["operation"] = operation
	}

	var buffer bytes.Buffer
	if encodeErr := json.NewEncoder(&buffer).Encode(data); encodeErr == nil {
		w.WriteHeader(statusCode)
		buffer.WriteTo(w)
	} else {
		//.. In the exceptional case where an error occurs while sending an error.
		logger.WriteError(encodeErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		var errorBytes, _ = json.Marshal(map[string]string{
			"encodeErr": encodeErr.Error(),
		})
		w.Write(errorBytes)
	}
}
