# Go-Inventory API

Sistem backend REST API untuk mengelolah manajemen inventaris.

## Teknologi yang digunakan

- **Go 1.23+**
- MySQL Database

## Features

- **Product** : Membuat produk baru, melihat seluruh produk, melihat produk berdasarkan product_id, menghapus produk, melakukan upload gambar produk dan men-download gambar produk.
- **Inventaris** : Melihat total stok dan lokasi produk seluruh inventaris produk, melihat total stok dan lokasi produk inventaris berdasarkan product_id, melakukan update stok dan lokasi produk.
- **Order** : Membuat pesanan, melihat pesanan berdasarkan order_id

## Setup Project

- Buat file .env yang menyimpan
  ```.env
  DB_USER = user
  DB_PASSWORD = password
  DB_HOST = localhost
  DB_PORT = 3306
  DB_NAME = "go-inventory"
  ```
- Create database SQL
  ```SQL
  create database `go-inventory`;
  ```
- Install semua dependency yang dibutuhkan terlebih dahulu
  ```bash
  go mod tidy
  ```
- Jalankan aplikasi. Tidak perlu create table karena project ini sudah menerapkan database migration.
  ```bash
  go run main.go
  ```
- Jalankan juga Postman test yang sudah tersedia di folder "./postman" untuk pengujian
