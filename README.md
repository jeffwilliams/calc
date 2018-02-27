# Building

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
    go get github.com/mna/pigeon github.com/cheekybits/genny
    go generate
    dep ensure
    go install


