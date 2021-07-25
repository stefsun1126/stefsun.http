FROM        golang:alpine
WORKDIR     /app
ENV         MODE=dev
COPY        . .
RUN         go mod download
RUN         go build -o app
CMD         ["./app"]


