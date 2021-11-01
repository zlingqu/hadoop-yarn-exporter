FROM golang:1.17.2-alpine3.14
WORKDIR /app
ADD . .
RUN apk add git \
    && go env -w GOPRIVATE=git-sa.nie.netease.com \
    && export GOPROXY=https://mirrors.aliyun.com/goproxy/ \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/ -v ./
FROM alpine:3.13.6
WORKDIR /app
COPY --from=0 /app/build/hadoop-yarn-exporter .
COPY --from=0 /app/krb5.conf .
COPY --from=0 /app/default.keytab .
CMD ["./app/hadoop-yarn-exporter"]

# docker build . -t ncr.nie.netease.com/ccgdc/hadoop-yarn-exporter:v20211101-1