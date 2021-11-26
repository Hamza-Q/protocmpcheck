.PHONY: usage
usage:
	echo "see Makefile"

ENV_SET = GOPATH=$(PWD) GO111MODULE=off  

.PHONY: install
install:
	$(ENV_SET) go get github.com/stretchr/testify google.golang.org/grpc google.golang.org/grpc/examples/helloworld/helloworld

.PHONY: test
test:
	cd ./src/testhelloworld && env $(ENV_SET) go test -v

.PHONY: print_env
print_env:
	@echo $(ENV_SET) 
