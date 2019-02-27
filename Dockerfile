FROM scratch

ADD ["bin/dnsns-bin", "/"]

CMD ["/dnsns-bin"]
