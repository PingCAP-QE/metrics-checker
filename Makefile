GOFILES := $(shell ls -r *.go)
MODULE_NAME := github.com/PingCAP-QE/metrics-checker
PACKAGES := go list ./...| grep -vE 'vendor' | grep 'github.com/PingCAP-QE/metrics-checker/'
PACKAGE_DIRECTORIES := $(PACKAGES) | sed 's|github.com/PingCAP-QE/metrics-checker||'

clean:
	rm -f metrics-checker

metrics-checker: go.mod $(GOFILES)
	go build -o $@

groupimports: install-goimports
	goimports -w -l -local github.com/PingCAP-QE/metrics-checker $$($(PACKAGE_DIRECTORIES))

install-goimports:
ifeq (,$(shell which goimports))
	@echo "installing goimports"
	go get golang.org/x/tools/cmd/goimports
endif

.PHONY: clean metrics-checker
