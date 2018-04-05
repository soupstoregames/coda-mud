FROM scratch

EXPOSE 50050/tcp

ADD ./bin/coda-world-linux-amd64 /coda-world

CMD ["/coda-world"]
