FROM scratch

COPY gopath/bin/cloudfunc /cloudfunc

ENTRYPOINT ["/cloudfunc"]