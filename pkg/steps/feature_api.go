package steps

import (
	"net/http/httptest"

	"github.com/cucumber/messages-go/v10"
)

type ApiFeature struct {
	resp *httptest.ResponseRecorder
}

func (a *ApiFeature) ResetResponse(*messages.Pickle) {
	a.resp = httptest.NewRecorder()
}

func (a *ApiFeature) SendRequest(method string, endpoint string, id int) error {
	return nil
}

func (a *ApiFeature) AssertResponseCode(code int) error {
	return nil
}

func (a *ApiFeature) AssertResponseBodyJSON(body *messages.PickleStepArgument_PickleDocString) error {
	return nil
}
