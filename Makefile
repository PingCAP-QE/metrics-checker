MODULE_NAME := github.com/PingCAP-QE/metrics-checker
PACKAGES := go list ./...| grep -vE 'vendor' | grep 'github.com/PingCAP-QE/metrics-checker/'
PACKAGE_DIRECTORIES := $(PACKAGES) | sed 's|github.com/PingCAP-QE/metrics-checker/||'


default: metrics-checker

clean:
	rm -rf bin/

metrics-checker:
	go build -o bin/metrics-checker cmd/metricchecker/*.go

groupimports: install-goimports
	goimports -w -l -local github.com/PingCAP-QE/metrics-checker $$($(PACKAGE_DIRECTORIES))

install-goimports:
ifeq (,$(shell which goimports))
	@echo "installing goimports"
	go get golang.org/x/tools/cmd/goimports
endif

.PHONY: clean metrics-checker
