package tests

import (
	"github.com/cucumber/godog"
	"github.com/nhood-org/nhood-engine-service-tests/pkg/steps"
)

func FeatureContext(s *godog.Suite) {
	api := steps.NewApiFeature()

	if mockEnabled() {
		go runMockAPIServer()
	}
	s.BeforeScenario(api.ResetResponse)

	s.Step(`^I send GET request to "([^"]*)" with ID (\d+)$`, api.SendRequest)
	s.Step(`^the response code should be (\d+)$`, api.AssertResponseCode)
	s.Step(`^the response should match json:$`, api.AssertResponseBodyJSON)
}
