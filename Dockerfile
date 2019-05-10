FROM scratch

ADD ["bin/dnsns.b", "/"]

CMD ["/dnsns.b"]
