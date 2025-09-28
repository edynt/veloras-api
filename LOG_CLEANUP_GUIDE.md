# Log Cleanup Guide

## Tổng quan

Hệ thống Veloras API đã được tích hợp tính năng tự động xóa file log cũ để quản lý dung lượng lưu trữ hiệu quả.

## Tính năng

### 1. Tự động xóa log cũ
- **Cron Job**: Chạy tự động mỗi ngày lúc 01:00 UTC
- **Tiêu chí**: Xóa file log cũ hơn 30 ngày
- **Log**: Ghi lại quá trình xóa trong log system

### 2. Xóa thủ công
- **API Endpoint**: `POST /api/v1/admin/cron/log-cleanup?days=30`
- **Web Interface**: Nút "🗑️ Cleanup Old Logs" trên trang log viewer
- **Tùy chỉnh**: Có thể thay đổi số ngày giữ lại

### 3. Hiển thị log dạng table
- **Format**: Bảng với các cột Level, Time, Caller, Message, Additional Data
- **Màu sắc**: Phân biệt theo level (ERROR=🔴, WARN=🟡, INFO=🔵, DEBUG=⚪)
- **Responsive**: Tương thích với mobile

## Cách sử dụng

### 1. Xem logs qua Web Interface
```bash
# Mở trình duyệt và truy cập
http://localhost:8080/web/logs.html
```

### 2. Xóa logs thủ công qua API
```bash
# Xóa logs cũ hơn 30 ngày
curl -X POST "http://localhost:8080/api/v1/admin/cron/log-cleanup?days=30"

# Xóa logs cũ hơn 7 ngày
curl -X POST "http://localhost:8080/api/v1/admin/cron/log-cleanup?days=7"
```

### 3. Kiểm tra scheduled jobs
```bash
curl "http://localhost:8080/api/v1/admin/cron/jobs"
```

### 4. Test log cleanup
```bash
# Chạy script test
./scripts/test-log-cleanup.sh
```

## Cấu hình

### 1. Thay đổi thời gian chạy cron job
Chỉnh sửa trong file `internal/cron/scheduler.go`:
```go
// Thay đổi từ "0 1 * * *" (01:00 UTC) thành thời gian khác
_, err = s.cron.AddFunc("0 2 * * *", s.logCleanupJob) // 02:00 UTC
```

### 2. Thay đổi số ngày mặc định
Chỉnh sửa trong file `internal/cron/scheduler.go`:
```go
// Thay đổi từ 30 thành số ngày khác
if err := s.logCleanupService.CleanupOldLogs(60); err != nil { // 60 ngày
```

### 3. Thay đổi thư mục log
Chỉnh sửa trong file `internal/cron/scheduler.go`:
```go
// Thay đổi đường dẫn thư mục log
logCleanupService := NewLogCleanupService("./custom/log/path")
```

## API Endpoints

### 1. Lấy danh sách file log
```
GET /api/v1/admin/logs/files
```

### 2. Lấy nội dung file log
```
GET /api/v1/admin/logs/content?filename=dev.001-2025-09-28.log&lines=100
```

### 3. Xóa logs thủ công
```
POST /api/v1/admin/cron/log-cleanup?days=30
```

### 4. Lấy thông tin scheduled jobs
```
GET /api/v1/admin/cron/jobs
```

## Monitoring

### 1. Log entries
Hệ thống sẽ ghi log khi:
- Bắt đầu cleanup job
- Xóa từng file log
- Hoàn thành cleanup job
- Lỗi trong quá trình cleanup

### 2. Log format
```json
{
  "level": "INFO",
  "time": "2025-09-28T01:00:00.000+0700",
  "caller": "cron/scheduler.go:110",
  "msg": "Starting log cleanup job",
  "scheduled_time": "2025-09-28 01:00:00"
}
```

## Troubleshooting

### 1. Cron job không chạy
- Kiểm tra `global.Config.Cron.Enabled` có được set thành `true`
- Kiểm tra log để xem có lỗi gì không
- Kiểm tra timezone của server

### 2. File log không được xóa
- Kiểm tra quyền ghi/xóa file trong thư mục log
- Kiểm tra file có đúng định dạng `.log` không
- Kiểm tra thời gian modification của file

### 3. Web interface không hiển thị
- Kiểm tra API endpoints có hoạt động không
- Kiểm tra CORS configuration
- Kiểm tra browser console để xem lỗi JavaScript

## Best Practices

1. **Backup**: Nên backup logs quan trọng trước khi xóa
2. **Monitoring**: Theo dõi log cleanup để đảm bảo hoạt động bình thường
3. **Testing**: Test cleanup trên môi trường dev trước khi deploy production
4. **Documentation**: Ghi lại các thay đổi cấu hình

## Security

- API endpoints được bảo vệ bởi middleware authentication
- Chỉ admin mới có thể chạy cleanup thủ công
- File path được validate để tránh directory traversal attacks
