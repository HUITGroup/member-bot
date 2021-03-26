FROM golang:latest AS build
RUN mkdir /memberbot
WORKDIR /memberbot
ADD . /memberbot
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o app .

FROM gcr.io/distroless/static-debian10
WORKDIR /root/
COPY --from=0 /memberbot/app .
CMD ["./app"]