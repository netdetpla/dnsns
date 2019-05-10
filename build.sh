#!/usr/bin/env bash
docker run -v /root/ndp/dnsns/bin/:/dnsns/bin/ -v /root/ndp/dnsns/dnsns/src/:/dnsns/src/ -e "CGO_ENABLED=0" dnsns-builder
