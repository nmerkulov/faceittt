FROM golang:1.15
WORKDIR /app
ADD . /app
RUN CGO_ENABLED=0 go build -o app main.go

FROM scratch
COPY --from=0 /app/app .
CMD ["./app"]