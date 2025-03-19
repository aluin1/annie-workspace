FROM golang:1.21-alpine

ARG APP_NAME
ENV APP_NAME=$APP_NAME
ENV TZ=Asia/Jakarta

# Install dependencies
RUN apk --no-cache add bash curl git gcc g++ tzdata dos2unix

WORKDIR /app

# Buat folder untuk tools
RUN mkdir -p /go/_tools

# Salin build script dan konfigurasi
COPY .docker/config/build.sh /go/_tools/build.sh
COPY .docker/config/reflex.conf /go/_tools/reflex.conf

# Konversi format Windows ke UNIX
RUN dos2unix /go/_tools/build.sh

# Berikan izin eksekusi pada script
RUN chmod +x /go/_tools/build.sh

# Download dan install reflex
RUN go install github.com/cespare/reflex@latest

# Tambahkan wait script
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

# Jalankan build.sh sebelum reflex
CMD ["/bin/sh", "/go/_tools/build.sh"]
