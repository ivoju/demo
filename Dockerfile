# GENERATE GO BINARY
FROM golang:1.13.3 as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/demo
COPY . ./

# download dependency
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /demo .
RUN cp -r env /env

# RUNNING GO BINARY
# Running go binary from compiler on the machine
FROM alpine:3.8

# SET TZ
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

# copy env from the host & copy go binary from the compiler
COPY --from=builder /env ./env
COPY --from=builder /demo ./

ENTRYPOINT ["/demo"]