###################################################
## BUILDER
###################################################

FROM golang:1.11-alpine AS builder

WORKDIR /go/src/github.com/dirty49374/tesseract-operator
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o ./tesseract-operator ./cmd/manager/main.go
RUN mkdir ./build/tmp && chmod 777 ./build/tmp

###################################################
## REAL
###################################################
FROM golang:1.11-alpine

WORKDIR /app
COPY --from=builder /go/src/github.com/dirty49374/tesseract-operator/tesseract-operator /app
COPY --from=builder /go/src/github.com/dirty49374/tesseract-operator/config/ /app/config/
#COPY --from=builder /go/src/github.com/dirty49374/tesseract-operator/secret/ /app/secret/

ENV WATCH_NAMESPACE=""

ENTRYPOINT [ "/app/tesseract-operator" ]
