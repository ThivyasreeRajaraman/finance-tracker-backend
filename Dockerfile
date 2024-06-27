FROM --platform=linux/amd64 golang:alpine as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM --platform=linux/amd64 nginx:alpine

WORKDIR /usr/share/nginx/html
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]