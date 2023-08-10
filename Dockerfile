# syntax=docker/dockerfile:1

FROM node:lts-alpine as node

WORKDIR /web

COPY web /web/

RUN npm install
RUN npm run build

FROM alpine:3.18 as golang
RUN apk add --no-cache --update go gcc g++

WORKDIR /app

COPY --from=node /web/build /app/web/build/.

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=1 GOOS=linux go build -o /personal-asset-tracker

FROM alpine:3.18

ARG PORT=8080
ENV PORT=${PORT:-8080}
ARG DATABASE_FILE
ENV DATABASE_FILE=${DATABASE_FILE:-database.db}

COPY --from=golang /personal-asset-tracker .

EXPOSE $PORT

CMD ["/personal-asset-tracker"]
