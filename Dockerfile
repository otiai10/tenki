FROM alpine:latest as build

RUN apk add go
RUN go install github.com/otiai10/tenki@latest

FROM alpine:latest as exec
RUN apk add tzdata
COPY --from=build /root/go/bin/tenki /bin

CMD ["tenki"]
