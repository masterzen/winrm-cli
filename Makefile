NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

TEST?=$(shell go list ./... | grep -v vendor)

# Get the current full sha from git
GITSHA:=$(shell git rev-parse HEAD)
# Get the current local branch name from git (if we can, this may be blank)
GITBRANCH:=$(shell git symbolic-ref --short HEAD 2>/dev/null)

all: deps
	@mkdir -p bin/
	@printf "$(OK_COLOR)==> Building$(NO_COLOR)\n"
	@GO15VENDOREXPERIMENT=1 go build -o $(GOPATH)/bin/winrm .

deps:
	@printf "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)\n"
	@printf "$(OK_COLOR)==> Installing godep and restoring dependencies$(NO_COLOR)\n"; \
		go get github.com/tools/godep; \
		godep restore;

clean:
	@rm -rf bin/ pkg/ src/

format:
	go fmt `go list ./... | grep -v vendor`

ci: deps
	@printf "$(OK_COLOR)==> Testing with Coveralls...$(NO_COLOR)\n"
	"$(CURDIR)/scripts/test.sh"

test: deps
	@printf "$(OK_COLOR)==> Testing...$(NO_COLOR)\n"
	@go test $(TEST) $(TESTARGS) -timeout=2m

updatedeps:
	@printf "$(ERROR_COLOR): winrm-cli deps are managed by godep.$(NO_COLOR)\n"

# This is used to add new dependencies to winrm-cli. If you are submitting a PR
# that includes new dependencies you will need to run this.
vendor:
	godep restore
	godep save

.PHONY: all clean deps format test updatedeps
