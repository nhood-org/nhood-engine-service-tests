package tests

import (
	"log"
	"net/http"
	"os"
)

const (
	testAgainstMockEnvVariable = "TEST_AGAINST_MOCK"
	testAgainstMockON          = "on"
	defaultPort                = ":8080"
)

func mockEnabled() bool {
	return os.Getenv(testAgainstMockEnvVariable) == testAgainstMockON
}

func runMockAPIServer() {
	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		response := ""
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	http.HandleFunc("/find", func(w http.ResponseWriter, r *http.Request) {
		size := r.URL.Query().Get("size")

		var status int
		body := ""

		switch size {
		default:
			status = http.StatusInternalServerError
		case "0":
			status = http.StatusBadRequest
		case "1":
			status = http.StatusOK
			body = `[{"id":"ID_0","key":["0.0","0.0","0.0"]}]`
		case "3":
			status = http.StatusOK
			body = `[{"id":"ID_0","key":["0.0","0.0","0.0"]},{"id":"ID_1","key":["0.0","0.0","1.0"]},{"id":"ID_5","key":["1.0","0.0","0.1"]}]`
		}

		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})

	log.Fatal(http.ListenAndServe(defaultPort, nil))
}
