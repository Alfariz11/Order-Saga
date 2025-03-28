# Rizki Alfariz Ramadhan/122140061/PWL

# Order Saga Microservices

Sistem ini menerapkan **Saga Pattern** dalam arsitektur microservices menggunakan bahasa pemrograman **Go (Golang)** dan komunikasi antar service menggunakan **gRPC**.

## 🧩 Arsitektur Layanan

Terdiri dari tiga layanan utama:

1. **Order Service**
   - Endpoint: `CreateOrder`, `CancelOrder`
   - Menyimpan status order (PENDING, COMPLETED, CANCELLED)

2. **Payment Service**
   - Endpoint: `ProcessPayment`, `RefundPayment`
   - Menyimpan status pembayaran (SUCCESS, FAILED, REFUNDED)

3. **Shipping Service**
   - Endpoint: `StartShipping`, `CancelShipping`
   - Menyimpan status pengiriman (PENDING, SHIPPED, CANCELLED)

Semua layanan diatur oleh **Saga Orchestrator** yang mengatur alur transaksi dan kompensasi bila terjadi kegagalan.

## 🔁 Alur Saga

1. `CreateOrder` dipanggil oleh orchestrator
2. Jika berhasil, lanjut ke `ProcessPayment`
3. Jika berhasil, lanjut ke `StartShipping`
4. Jika gagal di shipping:
   - `CancelShipping`
   - `RefundPayment`
   - `CancelOrder`

Jika gagal di payment:
- `CancelOrder`

## 🧪 Simulasi Pengujian

### ✅ Berhasil Semua:
- `CreateOrder("solon", "sepatu")`

### ❌ Shipping Gagal:
- `CreateOrder("solon", "fail-shipping")`
- Trigger kompensasi penuh

### ❌ Payment Gagal:
- `CreateOrder("solon", "fail-payment")`
- Trigger cancel order saja

### ❌ Input Tidak Valid:
- Jika nama atau item kosong, order tidak dibuat

## 🛠️ Cara Menjalankan

1. Jalankan semua service:
```bash
go run order-service/server/main.go
go run payment-service/server/main.go
go run shipping-service/server/main.go
```

2. Jalankan orchestrator:
```bash
go run saga-orchestrator/main.go
```

## 📁 Struktur Folder

```
OrderSaga/
├── proto/
├── gen/
├── order-service/
├── payment-service/
├── shipping-service/
├── saga-orchestrator/
├── go.mod
└── README.md
```

## 📦 Teknologi
- Go Modules
- gRPC
- Protobuf
- Saga Pattern
- Clean Architecture per service

---

Dibuat oleh: **Alfariz11** ✨
