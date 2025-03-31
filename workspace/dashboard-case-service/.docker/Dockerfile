# Menggunakan image Nginx sebagai web server
FROM nginx:alpine

# Menghapus konfigurasi default Nginx
RUN rm -rf /etc/nginx/conf.d/default.conf

# Menambahkan timezone Asia/Jakarta
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

# Menentukan direktori kerja
WORKDIR /usr/share/nginx/html

# Menyalin file frontend ke dalam container
COPY ./html /usr/share/nginx/html

COPY .docker/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 8080 untuk akses HTTP
EXPOSE 8080

# Menjalankan Nginx
CMD ["nginx", "-g", "daemon off;"]
