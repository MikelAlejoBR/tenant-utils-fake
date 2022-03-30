package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = "12000"

func main() {
	// We don't mind if the environment variable is empty.
	host := os.Getenv("HOST")

	// In this case we do mind if it is empty: in that case, use the default value.
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	address := fmt.Sprintf("%s:%s", host, port)

	// Both handlers can be tied to the same handler, since the expected body and the output's structure is the same.
	http.HandleFunc("/internal/orgIds", translate)
	http.HandleFunc("/internal/ebsNumbers", translate)

	log.Printf(`Listening on "%s"...`, address)
	log.Fatalln(http.ListenAndServe(address, nil))
}

// translate expects an array of strings as the body of the incoming request. For each one of those strings, it
// generates a random uint number, and returns everything in the format that the tenant translator expects it:
// {
//    "incomingNumber": "translatedNumber",
//    "incomingNumber2": "translatedNumber2",
// }
func translate(w http.ResponseWriter, req *http.Request) {
	requestBody, readErr := ioutil.ReadAll(req.Body)

	if err := req.Body.Close(); err != nil {
		errorMsg := fmt.Sprintf(`Could not close response's body: %s`, err)
		sendResponse(w, http.StatusInternalServerError, errorBody(errorMsg))
		log.Fatalf(errorMsg)
		return
	}

	if readErr != nil {
		errorMsg := fmt.Sprintf(`Error reading the response's body: %s`, readErr)
		sendResponse(w, http.StatusInternalServerError, errorBody(errorMsg))
		log.Fatalf(errorMsg)
		return
	}

	var numbers []string
	if err := json.Unmarshal(requestBody, &numbers); err != nil {
		errorMsg := fmt.Sprintf(`Error unmarshalling the response's body: %s`, err)
		sendResponse(w, http.StatusInternalServerError, errorBody(errorMsg))
		log.Fatalf(errorMsg)
		return
	}

	var results = make(map[string]string, len(numbers))
	for _, number := range numbers {
		results[number] = strconv.FormatUint(rand.Uint64(), 10)
	}

	responseBody, err := json.Marshal(results)
	if err != nil {
		errorMsg := fmt.Sprintf(`Could not marshal response to JSON: %s`, err)
		sendResponse(w, http.StatusInternalServerError, errorBody(errorMsg))
		log.Fatalf(errorMsg)
		return
	}

	sendResponse(w, http.StatusOK, responseBody)
	log.Printf(`Translation performed and response sent to "%s"`, req.RemoteAddr)
}

// sendResponse sends an "application/json" response with the given status and the given body.
func sendResponse(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Fatalf(`Could not send response: %s`, err)
		return
	}
}

// errorBody returns a byte array with the {"error": "message"} structure.
func errorBody(message string) []byte {
	return []byte(
		fmt.Sprintf(`{"error": "%s"}`, message),
	)
}
