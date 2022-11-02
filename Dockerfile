# Builder
FROM golang:1.19-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add make

WORKDIR /app

COPY ./ ./
RUN make build

# Runner
FROM alpine:latest

EXPOSE 8000

COPY --from=builder /app/company_http /app/company_http

WORKDIR /app
CMD ["./company_http"]
