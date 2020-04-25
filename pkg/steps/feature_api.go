package steps

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	resp *http.Response
}

func NewApiFeature() *ApiFeature {
	host := os.Getenv(targetHostEnvVariable)
	if host == "" {
		fmt.Printf("Environment variable %s not set, defaulting to `%s`\n", targetHostEnvVariable, targetHostDefault)
		host = targetHostDefault
	}

	return &ApiFeature{
		host: host,
	}
}

func (a *ApiFeature) ResetResponse(*messages.Pickle) {
	a.resp = nil
}

func (a *ApiFeature) SendRequest(endpoint string, id int) error {
	path := fmt.Sprintf(pathPattern, a.host, endpoint, id)

	resp, err := http.Get(path)
	if err != nil {
		return errors.Wrapf(err, "could not send request to %s", path)
	}

	a.resp = resp

	return nil
}

func (a *ApiFeature) AssertResponseCode(code int) error {
	if a.resp.StatusCode != code {
		m := fmt.Sprintf("actual response code '%d' does not match expected '%d'", a.resp.StatusCode, code)
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
