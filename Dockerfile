FROM golang:alpine AS build
WORKDIR /go/src/github.com/cuigh/swirl/
ADD . .
#RUN dep ensure
RUN go build

FROM alpine:3.6
MAINTAINER cuigh <noname@live.com>
WORKDIR /app
COPY --from=build /go/src/github.com/cuigh/swirl/swirl .
COPY --from=build /go/src/github.com/cuigh/swirl/config ./config/
COPY --from=build /go/src/github.com/cuigh/swirl/assets ./assets/
COPY --from=build /go/src/github.com/cuigh/swirl/views ./views/
EXPOSE 8001
ENTRYPOINT ["/app/swirl"]
