FROM node:20.12-alpine3.19 as tailwind

WORKDIR /node

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY components/*.templ ./components/
COPY tailwind.config.js .

RUN npx tailwindcss -o ./main.css --minify 

FROM golang:1.22.1-bullseye as builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY components/*.templ ./components/
COPY main.go .

RUN templ generate
RUN go build \
  -ldflags="-extldflags -static" \
  -o server \
  main.go


FROM alpine:3.19 as deploy

EXPOSE 8080

WORKDIR /app

COPY migrations/ /app/migrations/
COPY static/ /app/static/
COPY --from=tailwind /node/main.css /static/styles/main.css
COPY --from=builder /app/server .

CMD ["./server"]


