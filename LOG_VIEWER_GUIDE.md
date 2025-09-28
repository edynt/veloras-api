# 📊 Log Viewer & Cron Job Management Guide

## Tổng quan

Hệ thống đã được cập nhật với các tính năng sau:
- **Log ghi theo ngày**: Log files được tạo theo format `dev.001-YYYY-MM-DD.log`
- **Web interface**: Giao diện web để xem logs dễ dàng
- **API endpoints**: Các API để quản lý cron jobs và xem logs
- **Manual cleanup**: Khả năng chạy cleanup thủ công để test

## 🚀 Cách sử dụng

### 1. Khởi động server
```bash
go run cmd/server/main.go
```

### 2. Truy cập Web Interface
Mở trình duyệt và truy cập: **http://localhost:8080**

Hoặc truy cập trực tiếp: **http://localhost:8080/web/logs.html**

### 3. Sử dụng API Endpoints

#### Cron Job Management
```bash
# Xem danh sách cron jobs đã lên lịch
curl -X GET http://localhost:8080/api/v1/admin/cron/jobs

# Chạy cleanup thủ công
curl -X POST http://localhost:8080/api/v1/admin/cron/cleanup
```

#### Log Management
```bash
# Xem danh sách log files
curl -X GET http://localhost:8080/api/v1/admin/logs/files

# Xem nội dung log file (10 dòng cuối)
curl -X GET "http://localhost:8080/api/v1/admin/logs/content?filename=dev.001-2024-01-15.log&lines=10"
```

### 4. Chạy test script
```bash
./scripts/test-cron-and-logs.sh
```

## 📁 Cấu trúc Log Files

### Format tên file
- **Cũ**: `dev.001.log`
- **Mới**: `dev.001-2024-01-15.log` (theo ngày)

### Vị trí lưu trữ
```
storage/logs/
├── dev.001-2024-01-15.log
├── dev.001-2024-01-16.log
├── dev.001-2024-01-17.log
└── ...
```

## 🌐 Web Interface Features

### Tính năng chính
- **File selector**: Chọn log file từ dropdown
- **Line control**: Chọn số dòng hiển thị (10-1000)
- **Auto-refresh**: Tự động cập nhật danh sách files mỗi 30 giây
- **Syntax highlighting**: Màu sắc khác nhau cho các level log
- **Responsive design**: Tương thích mobile

### Cách sử dụng Web Interface
1. Mở **http://localhost:8080**
2. Chọn log file từ dropdown
3. Điều chỉnh số dòng muốn xem
4. Click "Load Log" để xem nội dung
5. Sử dụng nút refresh để cập nhật danh sách files

## ⚙️ Cron Job Configuration

### Lịch chạy tự động
- **Daily cleanup**: 00:01 UTC hàng ngày
- **Weekly stats**: 01:00 UTC thứ 2 hàng tuần

### Cấu hình trong `config.yaml`
```yaml
cron:
  enabled: true
  cleanup_sessions: true
  cleanup_tokens: true
  cleanup_verifications: true
  cleanup_password_resets: true
```

## 🔧 API Endpoints

### Cron Jobs
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/cron/jobs` | Xem danh sách cron jobs |
| POST | `/api/v1/admin/cron/cleanup` | Chạy cleanup thủ công |

### Log Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/logs/files` | Xem danh sách log files |
| GET | `/api/v1/admin/logs/content` | Xem nội dung log file |

### Parameters cho `/admin/logs/content`
- `filename` (required): Tên file log
- `lines` (optional): Số dòng hiển thị (default: 100)

## 📊 Log Levels & Colors

Trong web interface, các log levels được hiển thị với màu sắc khác nhau:
- **ERROR/FATAL**: 🔴 Đỏ (`#f48771`)
- **WARN**: 🟡 Vàng (`#dcdcaa`)
- **INFO**: 🔵 Xanh (`#9cdcfe`)
- **DEBUG**: ⚪ Xám (`#808080`)

## 🛠️ Troubleshooting

### Server không chạy
```bash
# Kiểm tra port 8080 có bị chiếm không
lsof -i :8080

# Khởi động server
go run cmd/server/main.go
```

### Không thấy log files
```bash
# Kiểm tra thư mục logs
ls -la storage/logs/

# Tạo thư mục nếu chưa có
mkdir -p storage/logs
```

### Web interface không load
- Kiểm tra server có chạy không
- Kiểm tra firewall/antivirus
- Thử truy cập trực tiếp: `http://localhost:8080/web/logs.html`

## 📝 Log Rotation

### Tự động rotation
- Log files được tạo mới mỗi ngày
- Files cũ được giữ lại theo cấu hình `MaxAge` (28 ngày)
- Files được compress khi cần thiết

### Cấu hình rotation
```yaml
logger:
  max_size: 500      # MB
  max_backups: 3     # Số files backup
  max_age: 28        # Ngày
  compress: true     # Nén files cũ
```

## 🎯 Best Practices

1. **Monitor logs thường xuyên** qua web interface
2. **Test cron jobs** bằng manual cleanup endpoint
3. **Backup log files** quan trọng
4. **Set up alerts** cho ERROR/FATAL logs
5. **Regular cleanup** của log files cũ

## 🔗 Useful Links

- **Web Interface**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html (nếu có)
- **Log Directory**: `./storage/logs/`
- **Test Script**: `./scripts/test-cron-and-logs.sh`
