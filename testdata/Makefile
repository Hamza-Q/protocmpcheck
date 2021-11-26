.PHONY: usage
usage:
	@echo "see Makefile"

# Testdata src assumes GOPATH mode for compatibility with analysistest.
# go tool invocations for use with that code must be set with these environment
# variables.
TESTDATA_GO_ENV = GOPATH=$(PWD) GO111MODULE=off  

# Print the environment for manual debugging and exploration of the testdata
# testhelloworld code.
.PHONY: print_env
print_env:
	@echo $(TESTDATA_GO_ENV) 

# Install the known testhelloworld dependencies
.PHONY: update_test
update_test:
	$(TESTDATA_GO_ENV) go get github.com/stretchr/testify google.golang.org/grpc google.golang.org/grpc/examples/helloworld/helloworld

# Run "go test" using the example code; this is nice for checking it
# compiles and fails as expected.
.PHONY: test
test: update_test
test:
	cd ./src/testhelloworld && env $(TESTDATA_GO_ENV) go test -v

# Run the linter on the example code and collect cpu profiles.
.PHONY: run_example_cpu
run_example_cpu: PROFILE_OPTIONS = -cpuprofile $(PWD)/_build/cpu.prof
run_example_cpu: run_example

# Run the linter on the example code and collect mem profiles.
.PHONY: run_example_mem
run_example_mem: PROFILE_OPTIONS = -memprofile $(PWD)/_build/mem.prof
run_example_mem: run_example

# Run the linter on the example code.
.PHONY: run_example
run_example: START_WD = $(PWD)
run_example: GOBIN=$(PWD)/_build
run_example:
	@# must be installed in module mode
	@echo "[make] Building protocmpcheck from source..."
	@cd .. && env GO111MODULE=on GOBIN=$(GOBIN) go install ./cmd/protocmpcheck
	@echo "[make] Running protocmpcheck on testdata/src/helloworld..."
	@time env $(TESTDATA_GO_ENV) $(GOBIN)/protocmpcheck $(PROFILE_OPTIONS) $(START_WD)/src/testhelloworld

