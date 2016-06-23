TEST?=$$(go list ./... | grep -v '/vendor/')
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

run:
	godep go run $(CURDIR)/cmd/arukas-ship/main.go

build: clean vet
	godep go build -o $(CURDIR)/bin/arukas-ship $(CURDIR)/cmd/arukas-ship/main.go

docker-build: clean vet
	docker build -t yamamotofebc/arukas-ship .

docker-run: docker-build
	docker run -it --rm -e ARUKAS_JSON_API_TOKEN \
		-e ARUKAS_JSON_API_SECRET \
		-e ARUKAS_ENDPOINT \
		-e ARUKAS_INSTANCE \
		-e ARUKAS_MEMORY \
		-e ARUKAS_CMD \
		-e SHIP_TOKEN \
		-e SHIP_PORT \
		-p $(SHIP_PORT):$(SHIP_PORT) \
		yamamotofebc/arukas-ship

docker-open: clean vet
	open "http://localhost:$(SHIP_PORT)?token=$(SHIP_TOKEN)"

test: vet
	godep go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

vet: fmt
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: default test vet fmt fmtcheck
