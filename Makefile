SRC=$(wildcard *.go)
GEN_SRC=gen_calc.go gen_eval.go
UNMANAGED_DEPS=github.com/cheekybits/genny github.com/mna/pigeon
UNMANAGED_DEPS_FULL=$(foreach dep, $(UNMANAGED_DEPS), $(GOPATH)/src/$(dep))
MANAGED_DEPS=vendor/github.com/chzyer/readline vendor/github.com/spf13/pflag

$(GOPATH)/bin/calc: $(GEN_SRC) $(SRC) $(MANAGED_DEPS)
	go test
	go install

$(GEN_SRC): calc.peg eval.genny $(UNMANAGED_DEPS_FULL)
	go generate

$(UNMANAGED_DEPS_FULL):
	go get $(UNMANAGED_DEPS)

$(MANAGED_DEPS): Gopkg.toml Gopkg.lock
	dep ensure

.PHONY: test
test: 
	go test

# Build for various architectures
archs:
	GOARCH=386 GOOS=linux go build -o calc_i386
