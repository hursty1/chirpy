package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)





func cleanString(jsonString string) string {
	var bannedWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
	}

	words := strings.Split(jsonString, " ")
	for i, word := range words {
		for _, banned := range bannedWords {
			if strings.Contains(strings.ToLower(word), strings.ToLower(banned)){
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}



func responseWithError(rw http.ResponseWriter, code int, msg string) {
	type responseError struct {
		Error string `json:"error"`
	}
	respErr := responseError{
		Error: msg,
	}
	jsonData, err := json.Marshal(respErr)
	if err != nil {
		log.Printf("Something has gone wrong marshalling json body %s", err)
		rw.WriteHeader(400)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(jsonData)
}
func responseWithJSON(rw http.ResponseWriter, code int, payload interface{}) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		responseWithError(rw,400,"Error Marshaling JSON data.")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(jsonData)
}

func responseWithOutBody(rw http.ResponseWriter, code int) {
	rw.WriteHeader(code)
}