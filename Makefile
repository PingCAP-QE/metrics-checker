GOFILES := $(shell ls -r *.go)
PACKAGES := go list ./...| grep -vE 'vendor'
PACKAGE_DIRECTORIES := $(PACKAGES) | sed 's|ofey404/metrics-checker/||'

clean:
	rm -f metrics-checker

metrics-checker: go.mod $(GOFILES)
	go build -o $@

groupimports: install-goimports
	goimports -w -l -local ofey404/metrics-checker $$($(PACKAGE_DIRECTORIES))

install-goimports:
ifeq (,$(shell which goimports))
	@echo "installing goimports"
	go get golang.org/x/tools/cmd/goimports
endif

.PHONY: clean metrics-checker
