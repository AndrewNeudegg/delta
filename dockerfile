FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git upx openssh
# Create appuser.
RUN adduser -D -g '' appuser
RUN mkdir /workdir
WORKDIR /workdir
COPY . .
RUN ls -alh
# Fetch dependencies.
# Using go get.
RUN go get -d -v ./cmd/delta
# Build the binary.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags="-w -s" -o /go/bin/src ./cmd/delta
# When you are ready to deploy uncomment the following line.
# RUN ls -lh /go/bin/src && echo "This step may take some time..." && date && upx -v --ultra-brute /go/bin/src && date

############################
# STEP 2 build a small image
############################
# FROM scratch
FROM alpine
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=builder /go/bin/src /delta
# Use an unprivileged user.
USER appuser
# Expose the port we will use.
EXPOSE 8080

# RUN ls -alh /
ENTRYPOINT ["/delta"]