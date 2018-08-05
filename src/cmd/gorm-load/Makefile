.PHONY: run deps-export

deps-export:
	@cd ../.. && glide install --strip-vendor

run: deps-export
	@go run *.go

run-quick:
	@go run *.go