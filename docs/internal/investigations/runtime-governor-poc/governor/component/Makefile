APP_NAME = runtime-governor
APP_PATH = docs/internal/investigations/runtime-governor-poc/governor/component
ENTRYPOINT = cmd/main.go
BUILDPACK = eu.gcr.io/kyma-project/test-infra/buildpack-golang-toolbox:v20200423-1d9d6590
SCRIPTS_DIR = $(realpath $(shell pwd)/../../../../../../)/scripts

include $(SCRIPTS_DIR)/generic_make_go.mk

errcheck-local:
	 GO111MODULE="on" errcheck -blank -asserts -ignoregenerated ./...

test-local:
	GO111MODULE=on go test ./...

resolve-local: ;

ensure-local: ;

dep-status-local: ;

push-image: ;
