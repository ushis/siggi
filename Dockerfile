FROM alpine

COPY siggi /usr/local/bin/siggi

USER 1

EXPOSE 8080

CMD siggi -listen :8080
