FROM golang:latest as builder

WORKDIR /go/src/GitOperator

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep init && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -o gitoperator -a -installsuffix cgo .

FROM alpine:latest 

RUN apk --no-cache add ca-certificates

RUN mkdir /app

WORKDIR /app

COPY --from=builder /go/src/GitOperator/gitoperator .

COPY --from=builder /go/src/GitOperator/.env .

EXPOSE 3030

CMD ["./gitoperator"]



