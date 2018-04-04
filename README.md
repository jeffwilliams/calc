# Building

## On Linux

    make 

## On other systems

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
    go get github.com/mna/pigeon github.com/cheekybits/genny
    go generate
    dep ensure
    go install
    go test

# TODO

Add support for decoding IP addresses from hex.

unit conversions. i.e. kg to lbs and ounces
    t=kg_to_lbs(2.2); floor(t); lbs_to_oz( t-floor(t))

list support:
  - fix code so testcasei `TestCalc/wrong_list_len` passes
	- allow list mmultiplication with number
	- add functions to extract list items
	- make conv funcs use lists
