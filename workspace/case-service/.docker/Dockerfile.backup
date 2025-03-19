FROM golang:1.21-alpine

ARG APP_NAME
ENV APP_NAME=$APP_NAME
ENV TZ=Asia/Jakarta


# Install curl 
RUN apk --no-cache add bash \
	curl \
	git \
	gcc \
	g++ \
	tzdata

WORKDIR /app


RUN mkdir -p /go/_tools

ADD .docker/config/build.sh /go/_tools
ADD .docker/config/reflex.conf /go/_tools

RUN chmod +x /go/_tools/build.sh

# Download and install reflex
RUN go install github.com/cespare/reflex@latest

# Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

# CMD will run the wait script first and then reflex
CMD /wait && reflex -c /go/_tools/reflex.conf