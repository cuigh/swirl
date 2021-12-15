# ---- Build UI----
FROM node:alpine AS node
WORKDIR /app
COPY ui .
RUN yarn install
RUN yarn run build

# ---- Build Go----
FROM golang:1.17-alpine AS golang
WORKDIR /app
COPY --from=node /app/dist ui/dist
COPY . .
RUN apk update && apk add git
RUN CGO_ENABLED=0 go build -ldflags "-s -w"

# ---- Release ----
FROM alpine
LABEL maintainer="cuigh <noname@live.com>"
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=golang /app/swirl .
COPY --from=golang /app/config config/
EXPOSE 8001
ENTRYPOINT ["/app/swirl"]