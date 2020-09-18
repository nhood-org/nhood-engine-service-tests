package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/cucumber/messages-go/v10"
	"github.com/pkg/errors"
)

const (
	targetHostEnvVariable = "TEST_TARGET_HOST"
	targetHostDefault     = "localhost:8080"
	pathAddDataPattern    = "http://%s/data"
	pathFindPattern       = "http://%s/find?size=%d"
	contentType           = "application/json"
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

func (a *ApiFeature) AddDataSetToService(table *messages.PickleStepArgument_PickleTable) error {
	path := fmt.Sprintf(pathAddDataPattern, a.host)

	dataToAdd := NewDataSliceFrom(table)
	for _, d := range dataToAdd {
		if d.ID == "" || d.Key == nil {
			continue
		}

		body, err := json.Marshal(d)
		if err != nil {
			return errors.Wrapf(err, "could not marshall data of id '%s' to json", d.ID)
		}

		resp, err := http.Post(path, contentType, bytes.NewReader(body))
		if err != nil {
			return errors.Wrapf(err, "could not send data of id '%s'", d.ID)
		}

		expectedCode := 200
		if resp.StatusCode != expectedCode {
			m := fmt.Sprintf("actual response code '%d' does not match expected '%d'", resp.StatusCode, expectedCode)
			return errors.New(m)
		}
	}

	return nil
}

func (a *ApiFeature) SendFindRequest(size int, table *messages.PickleStepArgument_PickleTable) error {
	path := fmt.Sprintf(pathFindPattern, a.host, size)

	if len(table.GetRows()) != 1 {
		return errors.New("there must be exactly one body element defined")
	}

	d := NewDataFrom(table.GetRows()[0])
	if d.ID == "" || d.Key == nil {
		return errors.Errorf("body element %v is invalid", d)
	}

	body, err := json.Marshal(d)
	if err != nil {
		return errors.Wrapf(err, "could not marshall data of id '%s' to json", d.ID)
	}

	resp, err := http.Post(path, contentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrapf(err, "could not send data of id '%s'", d.ID)
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

func (a *ApiFeature) AssertResponseElements(table *messages.PickleStepArgument_PickleTable) error {
	bodyBytes, err := ioutil.ReadAll(a.resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body bytes")
	}

	var data []Data
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return errors.Wrap(err, "could unmarshall body bytes to slice of data")
	}

	expectedData := NewDataSliceFrom(table)
	if !reflect.DeepEqual(data, expectedData) {
		m := fmt.Sprintf("actual body '%v' does not match expected '%v'", data, expectedData)
		return errors.New(m)
	}

	return nil
}
