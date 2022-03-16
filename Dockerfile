FROM golang:1.17-stretch AS build 
ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN apt install git make openssl

WORKDIR /go/src/github.com/alfianabdi/mutate-ingress-for-imported-service
ADD . .
# RUN make test
RUN make app

FROM scratch
WORKDIR /app
COPY --from=build /go/src/github.com/alfianabdi/mutate-ingress-for-imported-service/mutate-ingress-for-imported-service .
CMD ["/app/mutate-ingress-for-imported-service"]