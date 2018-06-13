FROM scratch

EXPOSE 50050/tcp

ADD ./bin/coda-linux-amd64 /coda

CMD ["/coda"]
