package steps

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/cucumber/messages-go/v10"
	"github.com/pkg/errors"
)

const (
	targetHostEnvVariable = "TEST_TARGET_HOST"
	targetHostDefault     = "localhost:8080"
	pathPattern           = "http://%s%s%d"
)

type ApiFeature struct {
	host string
	port string
	resp *httptest.ResponseRecorder
}

func NewApiFeature() *ApiFeature {
	host := os.Getenv(targetHostEnvVariable)
	if host == "" {
		fmt.Printf("Environment variable %s not set, defaulting to `%s`\n", targetHostEnvVariable, targetHostDefault)
		host = targetHostDefault
	}

	return &ApiFeature{
		host: host,
		resp: httptest.NewRecorder(),
	}
}

func (a *ApiFeature) ResetResponse(*messages.Pickle) {
	a.resp = httptest.NewRecorder()
}

func (a *ApiFeature) SendRequest(method string, endpoint string, id int) error {
	path := fmt.Sprintf(pathPattern, a.host, endpoint, id)

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return errors.Wrapf(err, "could not send request to %s", path)
	}

	http.DefaultServeMux.ServeHTTP(a.resp, req)

	return nil
}

func (a *ApiFeature) AssertResponseCode(code int) error {
	if a.resp.Code != code {
		m := fmt.Sprintf("actual response code '%d' does not match expected '%d'", a.resp.Code, code)
		return errors.New(m)
	}

	return nil
}

func (a *ApiFeature) AssertResponseBodyJSON(body *messages.PickleStepArgument_PickleDocString) error {
	bodyBytes, err := ioutil.ReadAll(a.resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body bytes")
	}

	bodyString := string(bodyBytes)
	if body.Content != bodyString {
		m := fmt.Sprintf("actual body '%s' does not match expected '%s'", body.Content, bodyString)
		return errors.New(m)
	}

	return nil
}
