package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type StatusRequest struct {
	Status string `json:"status"` // on, off, cycle, toggle
}

func getPortStatus(port string) (string, error) {
	cmd := exec.Command("/sbin/uhubctl", "-l", port)
	out, err := cmd.Output()
	if err != nil {
		// Throw an error if something goes wrong with the command
		return "", errors.New("Error running uhubctl: " + err.Error())
	}

	// Parse the output
	regex := regexp.MustCompile(`Port ` + port + `: [0-9]{4} ([a-z]+)`)
	status := regex.FindStringSubmatch(string(out))[1]
	return status, nil
}

func setPortStatus(port string, status string) (string, error) {
	statusMap := map[string]string{
		"off":    "0",
		"on":     "1",
		"cycle":  "2",
		"toggle": "3",
	}

	// Throw an error if the status is not defined in the map
	if _, ok := statusMap[status]; !ok {
		return "", errors.New("Invalid status: " + status)
	}

	cmd := exec.Command("/sbin/uhubctl", "-l", port, "-a", statusMap[status])
	out, err := cmd.Output()

	if err != nil {
		// Throw an error if something goes wrong with the command
		return "", errors.New("Error running uhubctl: " + err.Error())
	}

	// Parse uhubctl output
	newStatusRegex := regexp.MustCompile(`New status for hub ` + port)
	statusRegex := regexp.MustCompile(`Port ` + port + `: [0-9]{4} ([a-z]+)`)

	if newStatusRegex.MatchString(string(out)) && statusRegex.MatchString(string(out)) {
		return statusRegex.FindStringSubmatch(string(out))[1], nil
	}

	return "", errors.New("Error parsing uhubctl output")
}

// Takes the http request and returns the status of the port as json
func getStatus(w http.ResponseWriter, r *http.Request) {
	port := chi.URLParam(r, "port")
	status, err := getPortStatus(port)

	if err != nil {
		// If there is an error, return a 500 and the error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"status": "` + status + `"}`))
}

// Takes the http request and the new status from JSON body and returns the new status of the port as json
func setStatus(w http.ResponseWriter, r *http.Request) {
	port := chi.URLParam(r, "port")
	// Parse the JSON body
	var statusRequest StatusRequest
	err := json.NewDecoder(r.Body).Decode(&statusRequest)

	if err != nil {
		// If there is an error, return a 500 and the error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	// Set the port status
	status, err := setPortStatus(port, statusRequest.Status)

	if err != nil {
		// If there is an error, return a 500 and the error message
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"status": "` + status + `"}`))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It works!"))
	})

	r.Get("/ports/{port}", getStatus)
	r.Post("/ports/{port}", setStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.ListenAndServe(":"+port, r)
}
