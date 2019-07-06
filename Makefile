PACKAGE=github.com/x-armory/go-error
VERSION=$(shell cat ./VERSION)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_TAG=$(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_DIRTY} \

ifneq (${GIT_TAG},)
IMAGE_TAG=${GIT_TAG}
LDFLAGS += -X ${PACKAGE}.gitTag=${GIT_TAG}
endif

ALPINE_LDFLAGS += -w -linkmode external -extldflags "-static"

default: upgrade

help:
	@echo 'Management commands for lv:'
	@echo
	@echo 'Usage:'
	@echo '    make upgrade         Upgrade dependencies.'
	@echo '    make clean           Clean the directory tree.'
	@echo '    make test            Run tests on a compiled project.'
	@echo

.PHONY: upgrade
upgrade:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

.PHONY: clean
clean:
	@test ! -e bin || rm -rf bin

.PHONY: test
test:
	 go test ./...