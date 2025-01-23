# Gunakan gambar dasar Golang Alpine
FROM golang:1.21-alpine

# Set ARGs dan ENVs
ARG APP_NAME
ENV APP_NAME=$APP_NAME
# ENV TZ=Asia/Jakarta
ENV DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1
# Mode invariant

# Install dependensi
RUN apk --no-cache add \
    bash \
    curl \
    git \
    gcc \
    g++ \
    tzdata

# Install dependensi
RUN apk --no-cache add \
    bash \
    curl \
    icu-libs \
    libintl \
    libcurl \
    libunwind \
    ca-certificates \
    krb5-libs \
    && apk add --no-cache --virtual .build-deps \
    g++ gcc make musl-dev \
    && curl -sSL https://github.com/PowerShell/PowerShell/releases/download/v7.2.8/powershell-7.2.8-linux-alpine-x64.tar.gz -o /tmp/pwsh.tar.gz \
    && mkdir -p /opt/microsoft/powershell/7 \
    && tar -xvf /tmp/pwsh.tar.gz -C /opt/microsoft/powershell/7 \
    && ln -s /opt/microsoft/powershell/7/pwsh /usr/bin/pwsh \
    && apk del .build-deps

# Debugging: Cek versi PowerShell
RUN pwsh --version || echo "PowerShell tidak terinstal dengan benar"

# Set zona waktu
# RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo "${TZ}" > /etc/timezone

# Set direktori kerja
WORKDIR /app

# Siapkan tools dan skrip
RUN mkdir -p /go/_tools

# Tambahkan tools dan skrip
ADD .docker/config/build.sh /go/_tools/build.sh
ADD .docker/config/reflex.conf /go/_tools/reflex.conf
COPY Extract-Report.sh /app/

# Beri izin eksekusi pada skrip
RUN chmod +x /go/_tools/build.sh
RUN chmod +x /app/Extract-Report.sh

# Debugging: Verifikasi file yang ditambahkan
RUN ls -l /go/_tools && ls -l /app

# Install Reflex (File watcher)
RUN go install github.com/cespare/reflex@latest

# Tambahkan skrip wait
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

# CMD untuk menjalankan skrip dan tools
CMD /wait && /bin/bash /app/Extract-Report.sh && reflex -c /go/_tools/reflex.conf