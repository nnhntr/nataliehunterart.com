FROM alpine
RUN apk --no-cache add ca-certificates && apk --no-cache add tzdata
WORKDIR /app/
COPY . .
CMD ["/app/bin/nnhntr"]
