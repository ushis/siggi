FROM golang:onbuild

USER 1

EXPOSE 8080

CMD app -listen :8080
