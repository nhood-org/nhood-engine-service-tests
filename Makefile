default: clean install-dependencies install-tools run

GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard cmd/*.go)

export GO111MODULE = on

ARTIFACT_NAME = nhood-engine-service-tests

.PHONY: clean
clean:
	@echo "Cleaning:"
	go clean ./...
	@echo "...done"

.PHONY: install-tools
install-tools:
	@echo "Installing tools:"
	export GOBIN=$(GOBIN) && \
		go get github.com/cucumber/godog/cmd/godog
	@echo "...done"

.PHONY: install-dependencies
install-dependencies:
	@echo "Installing dependencies:"
	go mod vendor
	@echo "...done"

.PHONY: validate
validate:
	@echo "Validating tests:"
	export TEST_AGAINST_MOCK=on && \
		./bin/godog features/api.feature
	@echo "...done"

.PHONY: run
run:
	@test $(TEST_TARGET_HOST) || ( echo "TEST_TARGET_HOST not set" & exit 1 )
	@echo "Running tests:"
	export TEST_AGAINST_MOCK=off && \
		./bin/godog features/api.feature
	@echo "...done"

.PHONY: build-docker
build-docker:
	@test $(ARTIFACT_NAME) || ( echo "ARTIFACT_NAME not set" & exit 1 )
	@echo "Building docker image:"
	docker build -t nhood-org/${ARTIFACT_NAME}:local .
	@echo "...done"

.PHONY: build-docker-ci
build-docker-ci:
	@test $(ARTIFACT_NAME) || ( echo "ARTIFACT_NAME not set" & exit 1 )
	@test $(CIRCLE_BRANCH) || ( echo "CIRCLE_BRANCH not set" & exit 2 )
	@echo "Building docker image [CI]:"
	docker build -t nhood-org/${ARTIFACT_NAME}:${CIRCLE_BRANCH} .
	@echo "...done"

.PHONY: release-ci
release-ci:
	@test $(GITHUB_USERNAME) || ( echo "GITHUB_USERNAME not set" & exit 1 )
	@test $(GITHUB_TOKEN) || ( echo "GITHUB_TOKEN not set" & exit 2 )
	@test $(GITHUB_EMAIL) || ( echo "GITHUB_EMAIL not set" & exit 3 )
	@test $(NEW_VERSION) || ( echo "NEW_VERSION not set" & exit 4 )
	@echo "Releasing application version [CI]:"
	git config --global user.email ${GITHUB_EMAIL} && \
	git config --global user.name ${GITHUB_USERNAME} && \
	git tag -a v${NEW_VERSION} -m "${NEW_VERSION}" && \
	git push --tags https://${GITHUB_TOKEN}@github.com/nhood-org/${ARTIFACT_NAME}.git master
	@echo "...done"

.PHONY: release-docker-ci
release-docker-ci: build-docker-ci
	@test $(ARTIFACT_NAME) || ( echo "ARTIFACT_NAME not set" & exit 1 )
	@test $(GITHUB_USERNAME) || ( echo "GITHUB_USERNAME not set" & exit 2 )
	@test $(GITHUB_TOKEN) || ( echo "GITHUB_TOKEN not set" & exit 3 )
	@test $(CIRCLE_BRANCH) || ( echo "CIRCLE_BRANCH not set" & exit 4 )
	@test $(VERSION) || ( echo "VERSION not set" & exit 5 )
	@echo "Releasing docker image [CI]:"
	docker login docker.pkg.github.com -u ${GITHUB_USERNAME} -p ${GITHUB_TOKEN} && \
    docker tag nhood-org/${ARTIFACT_NAME}:${CIRCLE_BRANCH} docker.pkg.github.com/nhood-org/repository/${ARTIFACT_NAME}:${VERSION} && \
    docker push docker.pkg.github.com/nhood-org/repository/${ARTIFACT_NAME}:${VERSION}
	@echo "...done"

.PHONY: trigger-circle-ci-release
trigger-circle-ci-release:
	@test $(ARTIFACT_NAME) || ( echo "ARTIFACT_NAME not set" & exit 1 )
	@test $(CIRCLE_CI_USER_TOKEN) || ( echo "CIRCLE_CI_USER_TOKEN not set" & exit 2 )
	@test $(NEW_VERSION) || ( echo "NEW_VERSION not set" & exit 3 )
	@echo "Triggering application release:"
	curl -u ${CIRCLE_CI_USER_TOKEN}: \
		-d build_parameters[CIRCLE_JOB]=release \
		-d build_parameters[VERSION]=${NEW_VERSION} \
		https://circleci.com/api/v1.1/project/github/nhood-org/${ARTIFACT_NAME}/tree/master
	@echo "...done"

.PHONY: trigger-circle-ci-docker-release
trigger-circle-ci-docker-release:
	@test $(ARTIFACT_NAME) || ( echo "ARTIFACT_NAME not set" & exit 1 )
	@test $(CIRCLE_CI_USER_TOKEN) || ( echo "CIRCLE_CI_USER_TOKEN not set" & exit 2 )
	@test $(NEW_VERSION) || ( echo "NEW_VERSION not set" & exit 3 )
	@echo "Triggering docker release:"
	curl -u ${CIRCLE_CI_USER_TOKEN}: \
        -d build_parameters[CIRCLE_JOB]=release-docker \
        -d build_parameters[VERSION]=${NEW_VERSION} \
        https://circleci.com/api/v1.1/project/github/nhood-org/${ARTIFACT_NAME}/tree/master
	@echo "...done"
