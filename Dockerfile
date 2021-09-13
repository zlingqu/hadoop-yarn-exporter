FROM alpine:3.12.1
WORKDIR /app
ADD build .
ADD krb5.conf .
ADD default.keytab .
CMD ["/hadoop-yarn-exporter"]