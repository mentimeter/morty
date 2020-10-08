# Specify the version of Go to use
FROM golang:1.15-alpine

# Copy all the files from the host into the container
WORKDIR /src
COPY . .

# Compile the action
RUN go build -o /bin/action

ENTRYPOINT ["/bin/action"]
