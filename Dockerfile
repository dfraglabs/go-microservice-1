FROM alpine

ENV PORT 80
EXPOSE 80

RUN apk update && apk --no-cache add ca-certificates && update-ca-certificates

COPY build/go-microservice-1-linux /

ENTRYPOINT ["/go-microservice-1-linux", "-d"]
