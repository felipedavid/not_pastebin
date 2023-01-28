FROM golang:1.19-alpine as builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o not_pastebin ./cmd/web

FROM alpine:latest as runner
RUN mkdir /app
COPY --from=builder /app/not_pastebin /app

CMD ["/app/not_pastebin"]