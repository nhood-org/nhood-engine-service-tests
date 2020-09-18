package tests

import (
	"github.com/cucumber/godog"
	"github.com/nhood-org/nhood-engine-service-tests/pkg/steps"
)

func FeatureContext(s *godog.ScenarioContext) {
	api := steps.NewApiFeature()

	if mockEnabled() {
		go runMockAPIServer()
	}
	s.BeforeScenario(api.ResetResponse)

	s.Step(`^I add a set of data to the service:$`, api.AddDataSetToService)
	s.Step(`^I send a find request with expected size (\d+) and body:$`, api.SendFindRequest)
	s.Step(`^the response code should be (\d+)$`, api.AssertResponseCode)
	s.Step(`^the response contains the following elements:$`, api.AssertResponseElements)
}
