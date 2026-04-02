# ArbScanner (Crypto Arbitrage Scanner)

ArbScanner adalah aplikasi *bot scanner* ringan berbasis command-line (CLI) yang ditulis menggunakan bahasa Go. Aplikasi ini dirancang untuk mendeteksi peluang arbitrase (*arbitrage*) harga aset kripto secara *real-time* dengan membandingkan harga di berbagai bursa (*exchange*).

Saat ini, ArbScanner difokuskan untuk memantau pasangan mata uang **BTC-IDR** antara dua *exchange*:
1. **Indodax**
2. **Tokocrypto** (melalui jalur API Binance)

## Fitur Utama

- **Real-Time Data Streaming**: Mengambil data harga (*ticker*) secara berkala dari API *exchange*.
- **Arbitrage Detection Engine**: Secara otomatis membandingkan harga dari berbagai sumber dan akan memberikan notifikasi log jika mendeteksi perbedaan harga (*spread*) lebih dari 0.1%.
- **Concurrency Setup**: Menggunakan *Goroutines* dan *Channels* khas Go untuk menjalankan setiap *provider* secara paralel tanpa saling memblokir (*non-blocking*).
- **Clean Architecture**: Terbagi dalam *layer* `domain`, `adaptor`, dan `engine` agar kode mudah di-*maintenance* dan ditambahkan *exchange* baru.

## Prasyarat (Prerequisites)

Sebelum menjalankan aplikasi ini, pastikan sistem kamu memiliki:

1. **Go Environment**: Versi Go 1.25.5 atau yang lebih baru.
2. **VPN / Cloudflare WARP**: Karena API Binance sering kali dicegat oleh ISP Indonesia (Internet Positif), kamu diwajibkan untuk menyalakan VPN (seperti Cloudflare WARP / `warp-svc`) saat menjalankan aplikasi ini agar koneksi ke `api.binance.com` tidak *error*.

## Cara Menjalankan (How to Run)

1. Clone repositori ini ke komputer lokal kamu.
2. Pastikan VPN atau Cloudflare WARP kamu sudah dalam status **Connected**.
3. Buka terminal dan arahkan ke *root directory project* ini.
4. Jalankan perintah berikut untuk memulai *scanner*:

   ```bash
   go run cmd/app/main.go
   ```

## Struktur Direktori

Project ini menggunakan pendekatan Hexagonal/Clean Architecture:

    `/cmd/app`: Berisi `main.go` yang merupakan titik masuk (entry point) aplikasi. Di sinilah semua dependency disatukan.

    `/internal/domain`: Berisi kontrak antarmuka (`contract.go`) dan struktur data inti (`entity.go`) seperti MarketTicker.

    `/internal/adaptor`: Berisi implementasi spesifik untuk mengambil data dari masing-masing exchange (Indodax, Tokocrypto, dan Mock untuk testing).

    `/internal/engine`: Berisi logika bisnis utama (`detector.go`) yang mendengarkan aliran data dan menghitung kalkulasi persentase arbitrase.

## Disklaimer

Aplikasi ini dibuat murni untuk tujuan edukasi dan pemantauan pasar. Aplikasi belum terhubung ke sistem auto-trading atau eksekusi pesanan. Segala risiko finansial yang timbul dari keputusan perdagangan (trading) berdasarkan log aplikasi ini adalah tanggung jawab pengguna.

