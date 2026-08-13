# syntax=docker/dockerfile:1

FROM golang:1.26.4-alpine AS application
WORKDIR /app
RUN apk add --no-cache nodejs npm git build-base

# Copy root dependency tracking files if available, and app source code
COPY pkg/v1/win/systray/app/ ./pkg/v1/win/systray/app/

# Build frontend assets if required
WORKDIR /app/pkg/v1/win/systray/app/frontend
RUN npm install && npm run build

# Build Windows x64 (amd64) Wails Binary
WORKDIR /app/pkg/v1/win/systray/app
RUN GOOS=windows GOARCH=amd64 go build -tags desktop,production -ldflags="-s -w -H=windowsgui" -o app_amd64.exe .

# Build Windows x32 (386) Wails Binary
WORKDIR /app/pkg/v1/win/systray/app
RUN GOOS=windows GOARCH=386 go build -tags desktop,production -ldflags="-s -w -H=windowsgui" -o app_386.exe .

# Stage 2: Build the final system tray runner bundle (`./cmd/win`)
FROM golang:1.26.4-alpine AS bundler

WORKDIR /pkg/v1/win
COPY ./pkg/v1/win ./

WORKDIR /app/cmd/win
COPY cmd/win/ ./
RUN go mod edit -replace github.com/fehmicorp/go=/pkg || true
ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org,direct \
    GOCACHE=/root/.cache/go-build \
    GOMODCACHE=/root/go/pkg/mod

RUN go mod download || go mod download
RUN mkdir -p 64 86
COPY --from=application /app/pkg/v1/win/systray/app/app_amd64.exe ./64/app.exe
COPY --from=application /app/pkg/v1/win/systray/app/app_386.exe ./86/app.exe
RUN GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o ./64/connector.exe .
RUN GOOS=windows GOARCH=386 go build -ldflags="-s -w -H=windowsgui" -o ./86/connector.exe .

# Step 3 Create Installer
FROM ubuntu:latest AS installer
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /installer

RUN apt-get update && apt-get install -y --no-install-recommends \
    zip \
    nsis \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p 64 86

# Copy from bundler output path
COPY --from=bundler /app/cmd/win/64/ ./64/
COPY --from=bundler /app/cmd/win/86/ ./86/

# --- NSIS INSTALLER GENERATION ---
COPY cmd/installer.nsi ./installer.nsi
COPY cmd/win/assets/icon.ico ./icon.ico

RUN makensis -DSOURCE_DIR="/installer/64" -DOUTPUT_FILE="/installer/Setup-x64.exe" installer.nsi

RUN sed -i 's/\$PROGRAMFILES64/\$PROGRAMFILES/g' installer.nsi && \
    makensis -DSOURCE_DIR="/installer/86" -DOUTPUT_FILE="/installer/Setup-x86.exe" installer.nsi

# Final Step for output
FROM alpine:latest

WORKDIR /output
COPY --from=installer /installer/Setup-x64.exe ./
COPY --from=installer /installer/Setup-x86.exe ./

CMD ["echo", "Setup Installers created successfully at /output/"]