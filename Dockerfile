FROM scratch

ADD ["bin/dnsau-bin", "/"]

CMD ["/dnsau-bin"]
