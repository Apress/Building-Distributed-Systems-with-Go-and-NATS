FROM golang:1.23-bookworm AS build

WORKDIR /app

RUN ls
COPY chapter-4/go.mod ./
COPY chapter-4/go.sum ./
RUN go mod download

COPY ./chapter-4/customer ./
RUN pwd
RUN ls

RUN go build -o /server ./main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=build /server /server

USER nonroot:nonroot

ENTRYPOINT ["/server"]