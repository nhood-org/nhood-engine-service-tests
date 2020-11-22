package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"

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

func (a *ApiFeature) AddDataToService(table *messages.PickleStepArgument_PickleTable) error {
	path := fmt.Sprintf(pathAddDataPattern, a.host)

	dataToAdd := NewDataSliceFrom(table)

	if len(dataToAdd) != 1 {
		return errors.Errorf("expected exactly one data to add, while got %d", len(dataToAdd))
	}
	d := dataToAdd[0]

	if d.Reference == "" || d.Key == nil {
		return errors.New("invalid data: reference and key must be provided")
	}

	body, err := json.Marshal(d)
	if err != nil {
		return errors.Wrapf(err, "could not marshall data with reference '%s' to json", d.Reference)
	}

	a.resp, err = http.Post(path, contentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrapf(err, "could not send data with reference '%s'", d.Reference)
	}

	return nil
}

func (a *ApiFeature) AddDataSetToService(table *messages.PickleStepArgument_PickleTable) error {
	path := fmt.Sprintf(pathAddDataPattern, a.host)

	dataToAdd := NewDataSliceFrom(table)
	for _, d := range dataToAdd {
		if d.Reference == "" || d.Key == nil {
			continue
		}

		body, err := json.Marshal(d)
		if err != nil {
			return errors.Wrapf(err, "could not marshall data with reference '%s' to json", d.Reference)
		}

		resp, err := http.Post(path, contentType, bytes.NewReader(body))
		if err != nil {
			return errors.Wrapf(err, "could not send data with reference '%s'", d.Reference)
		}

		expectedCode := 201
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
	if d.Reference == "" || d.Key == nil {
		return errors.Errorf("body element %v is invalid", d)
	}

	body, err := json.Marshal(d)
	if err != nil {
		return errors.Wrapf(err, "could not marshall data with reference '%s' to json", d.Reference)
	}

	resp, err := http.Post(path, contentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrapf(err, "could not send data with reference '%s'", d.Reference)
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

func (a *ApiFeature) AssertResponseHeader(name string, exp string) error {
	r, err := regexp.Compile(exp)
	if err != nil {
		return errors.Wrapf(err, "could not compile expression '%s'", exp)
	}

	h := a.resp.Header.Get(name)
	if h == "" {
		return errors.Errorf("response does contain header of name '%s'", name)
	}

	matches := r.MatchString(h)
	if !matches {
		return errors.Errorf("header value '%s' does match given expression '%s'", h, exp)
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

	data = trimUUIDs(data)
	expectedData := NewDataSliceFrom(table)
	if !reflect.DeepEqual(data, expectedData) {
		m := fmt.Sprintf("actual body '%v' does not match expected '%v'", data, expectedData)
		return errors.New(m)
	}

	return nil
}

func trimUUIDs(data []Data) []Data {
	result := make([]Data, len(data))
	for i, d := range data {
		result[i] = Data{
			UUID:      "",
			Key:       d.Key,
			Reference: d.Reference,
		}
	}
	return result
}
