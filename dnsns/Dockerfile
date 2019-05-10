FROM golang:latest

RUN mkdir -p $GOPATH/src/golang.org/x/ \
    && cd $GOPATH/src/golang.org/x/ \
    && git clone https://github.com/golang/net.git \
    && git clone https://github.com/golang/text.git \
    && go get github.com/miekg/dns \
    && go get github.com/go-sql-driver/mysql \
    && go get github.com/weppos/publicsuffix-go/publicsuffix \
    && go get github.com/sparrc/go-ping

WORKDIR /dnssecure/src

CMD ["go", "build", "-o", "/dnssecure/bin/dnssecure.b", "./"]
