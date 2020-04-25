package tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		idString := strings.TrimPrefix(r.URL.Path, "/")
		id, err := strconv.Atoi(idString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Invalid ID" + idString))
			return
		}

		response := fmt.Sprintf(`{"id":%d}`, id)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	log.Fatal(http.ListenAndServe(defaultPort, nil))
}
