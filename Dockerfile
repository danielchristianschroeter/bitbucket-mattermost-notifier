FROM golang:latest

ARG VERSION=dev

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN go build -o main -ldflags=-X=main.version=${VERSION} .

# Set the environment variables
ENV LISTEN_PORT=1337
ENV MATTERMOST_WEBHOOKURL=
ENV MATTERMOST_USERNAME=Bitbucket
ENV MATTERMOST_CHANNEL=
ENV MESSAGE_ICON_EMOJI=:exclamation:
ENV INFO_COLOR=3498db
ENV DECLINED_COLOR=e74c3c
ENV WARNING_COLOR=f1c40f
ENV SUCCESS_COLOR=2ecc71

# Expose the listen port
EXPOSE ${LISTEN_PORT}

# Run the binary
CMD ["./main"]