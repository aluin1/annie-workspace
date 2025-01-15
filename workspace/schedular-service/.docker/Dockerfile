FROM golang:1.21-alpine

ARG APP_NAME
ENV APP_NAME=$APP_NAME
ENV TZ=Asia/Jakarta

# Install dependencies
RUN apk --no-cache add bash \
    curl \
    git \
    gcc \
    g++ \
    tzdata

# Instalasi PowerShell dengan cara manual
RUN mkdir -p /var/microsoft/powershell/7 \
    && curl -sSL https://github.com/PowerShell/PowerShell/releases/download/v7.2.8/powershell-7.2.8-linux-x64.tar.gz -o /tmp/pwsh.tar.gz \
    && tar -xvf /tmp/pwsh.tar.gz -C /var/microsoft/powershell/7 \
    && ls -l /var/microsoft/powershell/7 \
    && ln -s /var/microsoft/powershell/7/pwsh /usr/bin/pwsh

# Set timezone
RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo "${TZ}" > /etc/timezone

# Set working directory
WORKDIR /app

# Copy tools and scripts
RUN mkdir -p /go/_tools
ADD .docker/config/build.sh /go/_tools
ADD .docker/config/reflex.conf /go/_tools
COPY Extract-Report.sh /app/

# Provide execution permissions to scripts
RUN chmod +x /go/_tools/build.sh
RUN chmod +x /app/Extract-Report.sh

# Download and install Reflex
RUN go install github.com/cespare/reflex@latest

# Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

# CMD will run the wait script first and then reflex
CMD /wait && /bin/bash /app/Extract-Report.sh && reflex -c /go/_tools/reflex.conf
