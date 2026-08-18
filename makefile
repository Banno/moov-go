ifneq (,$(wildcard ./secrets.env))
    include secrets.env
    export
endif

.PHONY: build
build:
	go build ./...

# Pin lint-project.sh to an immutable moov-io/infra commit and verify its
# SHA-256 checksum before executing it. To bump: pick the new commit SHA,
# run `curl -fsSL <url> | sha256sum` locally, review the script, then update
# both values together.
# TODO(maintainers): replace these placeholders with the real commit SHA of
# moov-io/infra and the SHA-256 of go/lint-project.sh at that commit.
# `make check` fails closed until they are populated.
LINT_PROJECT_COMMIT := 0000000000000000000000000000000000000000
LINT_PROJECT_SHA256 := 0000000000000000000000000000000000000000000000000000000000000000
LINT_PROJECT_URL := https://raw.githubusercontent.com/moov-io/infra/$(LINT_PROJECT_COMMIT)/go/lint-project.sh

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	go test ./...
else
	@wget -O lint-project.sh $(LINT_PROJECT_URL)
	@if command -v sha256sum >/dev/null 2>&1; then \
		echo "$(LINT_PROJECT_SHA256)  lint-project.sh" | sha256sum -c -; \
	else \
		echo "$(LINT_PROJECT_SHA256)  lint-project.sh" | shasum -a 256 -c -; \
	fi || { echo "lint-project.sh checksum mismatch, refusing to execute" >&2; rm -f lint-project.sh; exit 1; }
	@chmod +x ./lint-project.sh
	COVER_THRESHOLD=30.0 GOTEST_PKGS=./... ./lint-project.sh
endif

.PHONY: clean
clean:
	@rm -rf ./bin/ ./tmp/ coverage.txt misspell* staticcheck lint-project.sh

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...

cover-web:
	go tool cover -html=cover.out

examples-e2ee:
	go test ./examples/e2ee/...