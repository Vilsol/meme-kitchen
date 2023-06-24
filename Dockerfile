FROM golang:alpine3.18 as backend

RUN apk add --no-cache libwebp-dev gcc g++

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 go build -v -a -installsuffix cgo -o /go/bin/api .


FROM node:18-alpine as frontend

WORKDIR /app

RUN npm i -g pnpm

COPY frontend/package.json package.json
COPY frontend/pnpm-lock.yaml pnpm-lock.yaml

RUN pnpm i

COPY frontend/ .

ENV NODE_ENV=production

RUN pnpm run build


FROM alpine

RUN apk add --no-cache libwebp

COPY --from=backend /go/bin/api /api
COPY --from=frontend /app/build/ static/
COPY fonts/ fonts/

ENTRYPOINT ["/api"]