.PHONY: test test-integration test-all

test:
	cd backend && go test ./...

test-integration:
	sh deploy/scripts/smoke-test.sh

test-all:
	sh deploy/scripts/test-all.sh
