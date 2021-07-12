GO := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go

.PHONY: run fmt

testbed: install-tcctl
	tcctl testbed create -f testbed.yaml -r http://172.16.4.180:30007

destroy-testbed:
	tcctl testbed delete -r http://172.16.4.180:30007

run:
	KUBECONFIG=$(KUBECONFIG) HTTP_PROXY=socks5://172.16.4.180:30081 go run main.go

build:
	$(GO) build -o bin/example main.go

image:
	GOOS=linux GOARCH=amd64 make build
	docker build -t hub.pingcap.net/qa/example_tidbcluster .

fmt: groupimports
	go fmt ./...

groupimports: install-goimports
	goimports -w -l -local github.com/pingcap/test-infra .

install-goimports:
ifeq (, $(shell which goimports))
	@{ \
	set -e ;\
	TMP_DIR=$$(mktemp -d) ;\
	cd $$TMP_DIR ;\
	GO111MODULE=on go get golang.org/x/tools/cmd/goimports ;\
	rm -rf $$TMP_DIR ;\
	}
endif

install-tcctl:
ifeq (, $(shell which tcctl))
	GOPRIVATE=github.com/pingcap go install github.com/pingcap/test-infra/tools/tcctl@latest
endif
