FROM golang:alpine AS build
WORKDIR /go/src/github.com/cuigh/swirl/
ADD . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w"

FROM alpine:3.8
LABEL maintainer="cuigh <noname@live.com>"
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/cuigh/swirl/swirl .
COPY --from=build /go/src/github.com/cuigh/swirl/config ./config/
COPY --from=build /go/src/github.com/cuigh/swirl/assets ./assets/
COPY --from=build /go/src/github.com/cuigh/swirl/views ./views/
EXPOSE 8001
ENTRYPOINT ["/app/swirl"]
