.PHONY: test cov-all gocov cov install-deps

install-deps:
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html/cmd/gocov-html@latest

gen:
	go generate ./...

test:
	go test -v ./...

cov-all:
	go test ./... -coverprofile=coverage/coverage-all.out
	go tool cover -html=coverage/coverage-all.out -o coverage/coverage-all.html
	go tool cover -func=coverage/coverage-all.out | grep total

cov:
	@PACKAGES=$$(go list ./... | grep -v "/mocks$$"| grep -v "/cmd$$"); \
	go test -coverpkg=$$(echo $$PACKAGES | tr ' ' ',') -coverprofile=coverage/coverage.out $$PACKAGES; \
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

gocov:
	@PACKAGES=$$(go list ./... | grep -v '/mocks' | grep -v '/cmd'); \
	gocov test $$PACKAGES > ./coverage/coverage-gocov.json; \
	gocov-html ./coverage/coverage-gocov.json > ./coverage/coverage-gocov.html
