FROM scratch

EXPOSE 5555/tcp

ADD ./bin/coda-linux-amd64 /coda

CMD ["/coda"]
