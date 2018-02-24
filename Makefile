SRC=$(wildcard *.go)

$GOPATH/bin/calc: peg.go $(SRC)
	go install

peg.go: calc.peg
	go generate

