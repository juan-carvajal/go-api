package utils

import (
	"fmt"
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, errorMessage string, externalMessage string, err error, httpCode int) {
	log.Println(fmt.Sprintf(errorMessage, err.Error()))
	http.Error(w, externalMessage, httpCode)
}
