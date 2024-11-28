# Use a Debian-based Go image
FROM golang:1.23-bullseye

WORKDIR /web

COPY . .

# Install tzdata for time zone information
RUN apt-get update && apt-get install -y tzdata

RUN go build -o app .

EXPOSE 80

CMD ["./app"]
