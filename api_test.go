package tests

import (
	"github.com/cucumber/godog"
	"github.com/nhood-org/nhood-engine-service-tests/pkg/steps"
)

var running = false

func FeatureContext(s *godog.ScenarioContext) {
	api := steps.NewApiFeature()

	if mockEnabled() && !running {
		go runMockAPIServer()
		running = true
	}
	s.BeforeScenario(api.ResetResponse)

	s.Step(`^I add data to the service:$`, api.AddDataToService)
	s.Step(`^I add a set of data to the service:$`, api.AddDataSetToService)
	s.Step(`^I send a find request with expected size (\d+) and body:$`, api.SendFindRequest)
	s.Step(`^the response code should be (\d+)$`, api.AssertResponseCode)
	s.Step(`^the response has header "([^"]*)" matching regex "([^"]*)"$`, api.AssertResponseHeader)
	s.Step(`^the response contains the following elements:$`, api.AssertResponseElements)
}
