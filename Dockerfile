FROM node:18 as tailwind

WORKDIR /node

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY components .

RUN npx tailwindcss -o ./main.css --minify 

FROM golang:1.22.1 as build

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY components .
COPY main.go .

RUN templ generate
RUN go build -o server .

FROM alpine as deploy

EXPOSE 8080

WORKDIR /app

COPY static .
COPY --from=tailwind /node/main.css /static/styles/main.css
COPY --from=build /app/server .

CMD ["./server"]


