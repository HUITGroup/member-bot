FROM golang:1.21 AS build
RUN mkdir /memberbot
WORKDIR /memberbot
ADD . /memberbot
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o app .

FROM gcr.io/distroless/static-debian11
COPY --from=0 /memberbot/app /usr/local/bin/member-bot
ENTRYPOINT [ "member-bot" ]
