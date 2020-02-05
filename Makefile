NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

TEST?=$(shell go list ./...)

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
		go mod tidy;

clean:
	@rm -rf bin/ pkg/ src/

format:
	go fmt $(go list ./...)

ci: deps
	@printf "$(OK_COLOR)==> Testing with Coveralls...$(NO_COLOR)\n"
	"$(CURDIR)/scripts/test.sh"

test: deps
	@printf "$(OK_COLOR)==> Testing...$(NO_COLOR)\n"
	@go test $(TEST) $(TESTARGS) -timeout=2m

updatedeps:
	@printf "$(ERROR_COLOR): winrm-cli deps are managed by go mod.$(NO_COLOR)\n"

.PHONY: all clean deps format test updatedeps
