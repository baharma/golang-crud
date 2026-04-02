# Alur Aplikasi

Dokumen ini menjelaskan alur kerja project secara konseptual dan operasional, tanpa membahas detail implementasi kode baris per baris.

Fokus dokumen ini adalah menjawab pertanyaan:

- aplikasi ini terdiri dari apa saja
- komponen mana yang menjalankan tugas tertentu
- bagaimana urutan proses saat setup, migrasi, dan menjalankan server
- apa yang terjadi ketika request masuk
- kapan `web`, `migrate`, dan `worker` dipakai
- apa saja batasan sistem saat ini

## Tujuan Project

Project ini disiapkan sebagai backend CRUD berbasis Go dengan penyimpanan data di MySQL.

Secara umum, tujuan arsitekturnya sederhana:

- ada satu proses untuk melayani HTTP request
- ada satu proses khusus untuk menyiapkan schema database
- ada ruang untuk proses background terpisah jika nanti dibutuhkan

Dengan pemisahan ini, tanggung jawab tiap proses menjadi lebih jelas dan tidak saling bercampur.

## Gambaran Besar Sistem

Komponen utama di project ini ada tiga:

- `web`
- `migrate`
- `worker`

Peran masing-masing:

- `web` bertugas menerima request dari client dan memberikan response
- `migrate` bertugas menyiapkan database dan tabel
- `worker` disiapkan untuk pekerjaan latar belakang yang tidak cocok dijalankan di request HTTP

Relasi sederhananya seperti ini:

```text
Client
  |
  v
Web Server
  |
  v
MySQL

Migrate Command
  |
  v
MySQL

Worker
  |
  v
MySQL / job source lain
```

## Komponen Sistem

### 1. Web

`Web` adalah proses utama aplikasi.

Tanggung jawabnya:

- membaca konfigurasi
- memastikan koneksi ke database tersedia
- membuka port HTTP
- menerima request dari client
- menjalankan logika aplikasi
- mengirim response ke client

Saat ini, endpoint yang sudah pasti tersedia adalah health check. Itu berarti proses web saat ini sudah cukup untuk menunjukkan aplikasi hidup, tetapi belum lengkap sebagai CRUD penuh.

### 2. Migrate

`Migrate` adalah command sekali jalan.

Artinya, proses ini tidak dirancang untuk terus hidup seperti server. Ia dijalankan ketika dibutuhkan, lalu selesai.

Tanggung jawabnya:

- membaca konfigurasi database
- konek ke server MySQL
- memastikan database aplikasi ada
- membuat database jika belum ada
- membuat atau menyesuaikan tabel sesuai model yang dipakai aplikasi

Alasan command ini dipisahkan dari web:

- startup web jadi lebih bersih
- migrasi bisa dikontrol manual
- cocok untuk pipeline deployment
- risiko perubahan schema saat server boot jadi lebih kecil

### 3. Worker

`Worker` adalah proses background.

Berbeda dengan `migrate`, worker biasanya bukan command sekali jalan. Worker biasanya hidup terus dan menunggu pekerjaan.

Contoh peran worker:

- memproses antrian pekerjaan
- mengirim email
- membuat laporan berkala
- sinkronisasi data
- menjalankan tugas terjadwal

Di project ini, worker belum dipakai sungguhan. Jadi secara arsitektur ia masih berupa slot yang sudah disiapkan, bukan komponen aktif.

## Aktor yang Terlibat

Dalam alur sistem ini ada beberapa aktor:

- developer
- aplikasi `migrate`
- aplikasi `web`
- MySQL
- client atau consumer API

Peran tiap aktor:

- developer menyiapkan environment dan menjalankan command
- `migrate` menyiapkan database
- `web` melayani trafik aplikasi
- MySQL menyimpan data
- client mengakses API

## Alur Setup Awal

Alur setup awal adalah proses menyiapkan semua kebutuhan agar aplikasi bisa dijalankan.

Urutannya seperti ini:

1. developer menyiapkan Go
2. developer menyiapkan MySQL
3. developer menyiapkan file `.env`
4. developer menjalankan migrasi
5. developer menjalankan web server
6. developer menguji endpoint health

Penjelasan per langkah:

### 1. Menyiapkan Go

Go diperlukan untuk:

- build aplikasi
- menjalankan command
- mengelola dependency
- menjalankan test

Kalau Go belum tersedia, maka aplikasi tidak bisa dijalankan sama sekali.

### 2. Menyiapkan MySQL

MySQL adalah penyimpanan utama data.

Yang harus tersedia:

- service MySQL aktif
- host dan port bisa dijangkau
- user database valid
- password valid

Tanpa MySQL, proses `migrate` dan `web` tidak akan bisa berjalan normal.

### 3. Menyiapkan `.env`

File `.env` adalah sumber konfigurasi utama saat ini.

Nilai penting yang perlu tersedia:

- port aplikasi
- host database
- port database
- nama database
- user database
- password database

Kalau konfigurasi salah:

- migrate bisa gagal
- web bisa gagal start
- aplikasi bisa mengarah ke database yang salah

### 4. Menjalankan Migrasi

Tahap ini menyiapkan database aplikasi.

Tujuannya:

- memastikan nama database aplikasi tersedia
- memastikan tabel yang dibutuhkan sudah ada

Tanpa migrasi, web server mungkin bisa hidup jika database sudah siap, tetapi struktur tabel belum tentu ada.

### 5. Menjalankan Web Server

Tahap ini mengaktifkan aplikasi yang menerima request.

Kalau proses ini berhasil:

- port aplikasi akan terbuka
- endpoint health bisa diakses
- aplikasi siap menerima trafik

### 6. Menguji Health Check

Ini adalah verifikasi paling dasar.

Jika health check sukses, artinya:

- proses web hidup
- port aplikasi terbuka
- routing dasar aktif

Health check tidak selalu berarti semua fitur bisnis sudah lengkap, tetapi cukup untuk membuktikan server siap berjalan.

## Alur Konfigurasi

Konfigurasi adalah pintu masuk semua proses.

Secara operasional, alurnya seperti ini:

1. aplikasi membaca `.env`
2. nilai konfigurasi dimuat ke memori
3. nilai itu dipakai untuk menentukan:
   - port server
   - target database
   - kredensial database

Dampak konfigurasi:

- konfigurasi port menentukan aplikasi didengar di port mana
- konfigurasi DB menentukan aplikasi bicara ke MySQL yang mana
- konfigurasi salah akan memicu error di startup

Karena itu, `.env` harus dianggap sebagai syarat wajib sebelum menjalankan command apa pun.

## Alur Migrasi

Alur migrasi sekarang adalah salah satu bagian paling penting di project ini.

Secara konsep, urutannya seperti ini:

1. command `migrate` dijalankan
2. command membaca konfigurasi database
3. command mencoba konek ke server MySQL
4. command memastikan database aplikasi ada
5. jika database belum ada, database dibuat
6. command konek ke database aplikasi
7. command membuat atau menyesuaikan tabel
8. command selesai

Diagram sederhana:

```text
Developer
  |
  v
Run migrate
  |
  v
Read .env
  |
  v
Connect to MySQL server
  |
  v
Ensure database exists
  |
  v
Connect to application database
  |
  v
Create or update tables
  |
  v
Done
```

Hasil yang diharapkan setelah migrasi:

- database aplikasi tersedia
- tabel utama tersedia
- web server bisa start tanpa gagal karena schema belum ada

## Alur Menjalankan Web

Setelah migrasi selesai, proses berikutnya adalah menjalankan `web`.

Alurnya:

1. command `web` dijalankan
2. aplikasi membaca konfigurasi
3. aplikasi mencoba konek ke database
4. jika koneksi berhasil, server HTTP dinyalakan
5. server mulai mendengarkan request di port yang ditentukan
6. client bisa mulai memanggil endpoint

Hal penting di alur ini:

- `web` tidak lagi bertugas membuat database
- `web` tidak lagi bertugas menjalankan migrasi schema
- `web` hanya memastikan database yang sudah disiapkan memang bisa dipakai

Keuntungan pendekatan ini:

- flow startup lebih jelas
- masalah migrasi terpisah dari masalah trafik HTTP
- debugging lebih mudah

## Alur Request dari Client

Walaupun endpoint CRUD belum dibuat, pola umumnya tetap sama.

Alur request HTTP secara konseptual:

1. client mengirim request ke aplikasi
2. server menerima request di port aplikasi
3. router menentukan handler yang tepat
4. aplikasi menjalankan logika yang dibutuhkan
5. jika perlu, aplikasi membaca atau menulis data ke MySQL
6. aplikasi membentuk response
7. response dikembalikan ke client

Untuk endpoint health, alurnya lebih pendek:

1. client memanggil `/health`
2. server menerima request
3. handler health menyiapkan response sederhana
4. response sukses dikirim ke client

Fungsi health check:

- memastikan server hidup
- memastikan routing dasar aktif
- memudahkan pengecekan lokal, deployment, dan monitoring

## Alur Kegagalan yang Mungkin Terjadi

Di luar kode, ini adalah skenario gagal yang paling umum.

### 1. MySQL tidak aktif

Dampaknya:

- `migrate` gagal konek
- `web` gagal start

Gejalanya:

- koneksi ditolak
- timeout
- server tidak bisa dijangkau

### 2. Kredensial database salah

Dampaknya:

- aplikasi tidak bisa autentikasi ke MySQL

Gejalanya:

- access denied
- authentication failed

### 3. Nama database salah

Dampaknya:

- `web` bisa gagal connect ke database target
- `migrate` bisa membuat database yang tidak sesuai harapan

### 4. Port aplikasi bentrok

Dampaknya:

- `web` gagal bind ke port

Gejalanya:

- address already in use

### 5. Struktur schema belum siap

Dampaknya:

- fitur yang membutuhkan tabel tertentu bisa gagal

Saat ini risiko ini sudah dikurangi dengan adanya command migrasi khusus.

## Alur Worker Secara Konseptual

Walaupun worker belum diimplementasikan, penting untuk memahami kapan worker dibutuhkan.

Worker dipakai ketika ada tugas yang:

- tidak perlu selesai dalam request yang sama
- butuh waktu lama
- lebih aman dijalankan terpisah
- sebaiknya diproses di latar belakang

Contoh alurnya nanti jika worker dipakai:

1. client melakukan aksi tertentu ke API
2. web menerima request
3. web menyimpan data utama
4. web mendaftarkan pekerjaan background
5. web langsung memberi response ke client
6. worker mengambil pekerjaan itu
7. worker memprosesnya
8. hasil dipersist ke database atau dikirim ke sistem lain

Keuntungan worker:

- response API lebih cepat
- pekerjaan berat tidak membebani request utama
- proses bisa diulang jika gagal

Karena worker belum aktif, alur itu saat ini masih belum berjalan di project.

## Alur Data

Kalau dilihat dari sisi data, project ini masih punya alur sederhana:

1. data masuk dari request
2. aplikasi memvalidasi dan memproses data
3. data disimpan atau dibaca dari MySQL
4. hasil dikembalikan ke client

Saat ini tabel yang sudah disiapkan mewakili:

- kategori
- buku

Artinya fondasi data sudah disiapkan untuk domain sederhana antara kategori dan buku.

## Alur Pengembangan Harian

Untuk developer, alur kerja paling sehat saat ini adalah:

1. update `.env` bila perlu
2. pastikan MySQL aktif
3. jalankan migrasi
4. jalankan test
5. jalankan web server
6. cek endpoint health
7. lanjut menambah endpoint atau logika bisnis

Kenapa urutannya seperti itu:

- migrasi memastikan schema siap
- test memastikan compile dan paket dasar aman
- web dijalankan setelah fondasi siap

## Alur Deployment Sederhana

Kalau nanti aplikasi mau dipindahkan ke server atau environment lain, flow sederhananya biasanya seperti ini:

1. siapkan environment variable
2. siapkan MySQL target
3. jalankan migrasi
4. jalankan proses web
5. pantau health check

Kalau worker nanti dipakai, urutannya menjadi:

1. siapkan environment
2. jalankan migrasi
3. jalankan web
4. jalankan worker
5. pantau health dan log

## Batasan Sistem Saat Ini

Secara non-code, sistem ini masih berada di tahap fondasi.

Yang sudah ada:

- alur setup dasar
- command migrasi
- koneksi database
- health endpoint
- pemisahan peran `web`, `migrate`, dan `worker`

Yang belum ada:

- endpoint CRUD penuh
- validasi alur bisnis end-to-end
- worker yang benar-benar bekerja
- mekanisme queue
- dokumentasi deployment rinci
- versioned migration yang lebih formal

## Kenapa Pemisahan Ini Penting

Memisahkan `web`, `migrate`, dan `worker` adalah keputusan arsitektur yang sehat karena:

- tiap proses punya tanggung jawab jelas
- kegagalan lebih mudah dilokalisasi
- deployment lebih fleksibel
- pengembangan jangka panjang lebih rapi

Kalau semua tanggung jawab dicampur ke dalam satu proses:

- startup jadi berat
- debugging jadi kabur
- perubahan schema bercampur dengan trafik runtime
- perluasan fitur jadi lebih sulit

## Ringkasan Inti

Kalau disederhanakan, alur besar project ini adalah:

1. siapkan konfigurasi
2. siapkan database
3. jalankan migrasi
4. jalankan web server
5. client mengakses API
6. worker disiapkan untuk kebutuhan background di tahap berikutnya

Jadi peran utamanya bisa dibaca begini:

- `migrate` menyiapkan tempat penyimpanan
- `web` melayani aplikasi
- `worker` nanti menangani pekerjaan latar belakang

Dokumen ini sengaja menjelaskan sistem dari sisi alur operasional, bukan dari sisi isi file kode. Jika dibutuhkan, tahap berikutnya yang paling masuk akal adalah membuat dokumen lanjutan untuk:

- alur CRUD buku dan kategori
- rancangan endpoint API
- rancangan worker bila nanti mulai dipakai
- alur deployment untuk server production
