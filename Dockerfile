FROM golang:alpine

COPY gopath/bin/cloudfunc /cloudfunc

ENTRYPOINT ["/cloudfunc"]