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