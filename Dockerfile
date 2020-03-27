FROM golang
WORKDIR /src
COPY hitwin.go config.json ./
RUN go build hitwin.go
FROM ubuntu
COPY --from=0 /src/. .
EXPOSE 8000
CMD ["./hitwin"]
