FROM node:18 as tailwind

WORKDIR /node

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY components .

RUN npx tailwindcss -o ./main.css --minify 

FROM golang:1.22.1 as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY components .
COPY main.go .

RUN go build . -o server

FROM alpine as deploy

WORKDIR /app

COPY static .
COPY --from=tailwind main.css /static/styles/main.css
COPY --from=build server .

CMD ["./server"]

