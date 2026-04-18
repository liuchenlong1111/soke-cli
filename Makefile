BINARY := soke-cli                                                                                                                                                                                                                           
VERSION := dev                                                                                                                                                               
LDFLAGS := -s -w -X main.Version=$(VERSION)                                                                                                                                                                                                    

.PHONY: build install test clean                                                                                                                                                                                                               

build:                                                                                                                                                                                                                                         
    go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: build                                                                                                                                                                                                                                 
    install -m755 $(BINARY) /usr/local/bin/$(BINARY)                                                                                                                                                                                         

test:           
    go test -v ./...                                                                                                                                                                                                                         
                                                                                                                                                                                                                                               
clean:
    rm -f $(BINARY)