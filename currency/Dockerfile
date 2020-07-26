# Use golang as the build container
FROM golang:1.14.6-alpine as build

RUN apk update && apk add --no-cache git protoc

# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# Copy the source code into the working directory
COPY . /go/src/github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition

# Set the working directory
WORKDIR /go/src/github.com/PacktPublishing/Building-Microservices-with-Go-Second-Edition/currency

# Build the protobufs
RUN go get github.com/golang/protobuf/protoc-gen-go && \
  protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /bin/currency .

# Use scratch as the run container
FROM scratch

# Import the user and group files from the builder.
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

# Copy the statically linked binary
COPY --from=build /bin/currency /bin/currency

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 9092

ENTRYPOINT [ "/bin/currency"]