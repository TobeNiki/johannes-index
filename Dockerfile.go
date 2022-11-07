FROM golang:1.18-alpine

RUN apk upgrade --update && \
    apk --no-cache add git

WORKDIR /backend
#GIN_MODE=release
ENV GIN_MODE=debug

# timezone install
RUN apk add --update --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata
#ENV TZ=Asia/Tokyo
COPY go.mod .
COPY go.sum . 
RUN go mod download

ENV PORT=${PORT}

#ENTRYPOINT [ "backend" ]
CMD [ "go", "run", "backend.go" ]