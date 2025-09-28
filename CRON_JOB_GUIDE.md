# Cron Job Cleanup Guide

## Tổng quan

Hệ thống cron job được thiết kế để tự động dọn dẹp các dữ liệu hết hạn trong database vào 00:01 UTC hằng ngày. Tính năng này giúp duy trì hiệu suất database và giải phóng không gian lưu trữ.

## Các loại dữ liệu được dọn dẹp

### 1. Sessions hết hạn
- **Bảng**: `sessions`
- **Điều kiện**: `expires_at < current_timestamp`
- **Mục đích**: Xóa các session đã hết hạn để bảo mật

### 2. Email Verification Codes hết hạn
- **Bảng**: `email_verifications`
- **Điều kiện**: `expires_at < current_timestamp`
- **Mục đích**: Xóa các mã xác thực email đã hết hạn

### 3. Password Reset Tokens hết hạn
- **Bảng**: `password_resets`
- **Điều kiện**: `expires_at < current_timestamp`
- **Mục đích**: Xóa các token reset password đã hết hạn

## Cấu hình

### 1. Config File (config.yaml)

```yaml
cron:
  enabled: true                    # Bật/tắt cron scheduler
  cleanup_sessions: true           # Dọn dẹp sessions hết hạn
  cleanup_tokens: true             # Dọn dẹp tokens hết hạn
  cleanup_verifications: true      # Dọn dẹp email verifications hết hạn
  cleanup_password_resets: true    # Dọn dẹp password resets hết hạn
```

### 2. Lịch trình chạy

- **Daily Cleanup**: 00:01 UTC hằng ngày
- **Weekly Stats**: 01:00 UTC thứ 2 hằng tuần (báo cáo thống kê)

## API Endpoints

### 1. Xem thông tin scheduled jobs

```http
GET /api/v1/admin/cron/jobs
```

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "schedule": "Daily at 00:01 UTC",
      "next_run": "2024-01-16T00:01:00Z",
      "last_run": "2024-01-15T00:01:00Z"
    }
  ],
  "error": false
}
```

### 2. Chạy cleanup thủ công

```http
POST /api/v1/admin/cron/cleanup
```

**Response:**
```json
{
  "code": 200,
  "message": "Success",
  "data": "Cleanup completed successfully",
  "error": false
}
```

## Monitoring và Logging

### 1. Logs

Hệ thống sẽ log các sự kiện quan trọng:

```
{"level":"INFO","time":"2024-01-15T00:01:00Z","msg":"Starting daily cleanup job"}
{"level":"INFO","time":"2024-01-15T00:01:05Z","msg":"Cleanup expired sessions completed","deleted_count":15}
{"level":"INFO","time":"2024-01-15T00:01:06Z","msg":"Cleanup expired email verifications completed","deleted_count":8}
{"level":"INFO","time":"2024-01-15T00:01:07Z","msg":"Cleanup expired password resets completed","deleted_count":3}
{"level":"INFO","time":"2024-01-15T00:01:10Z","msg":"Daily cleanup job completed successfully"}
```

### 2. Error Handling

Nếu có lỗi xảy ra:

```
{"level":"ERROR","time":"2024-01-15T00:01:05Z","msg":"Failed to cleanup expired sessions","error":"connection timeout"}
{"level":"ERROR","time":"2024-01-15T00:01:10Z","msg":"Daily cleanup job failed","error":"cleanup completed with 1 errors"}
```

## Cách sử dụng

### 1. Khởi động server

```bash
go run cmd/server/main.go
```

Cron scheduler sẽ tự động khởi động cùng với server.

### 2. Kiểm tra trạng thái

```bash
# Xem scheduled jobs
curl -X GET http://localhost:8080/api/v1/admin/cron/jobs

# Chạy cleanup thủ công
curl -X POST http://localhost:8080/api/v1/admin/cron/cleanup
```

### 3. Monitor logs

```bash
# Xem logs real-time
tail -f storage/logs/dev.001.log | grep -i cleanup
```

## Tùy chỉnh

### 1. Thay đổi lịch trình

Trong file `internal/cron/scheduler.go`, bạn có thể thay đổi cron expression:

```go
// Thay đổi từ "1 0 * * *" (00:01 UTC) thành "0 2 * * *" (02:00 UTC)
_, err := s.cron.AddFunc("0 2 * * *", s.dailyCleanupJob)
```

### 2. Thêm cleanup mới

1. Thêm query vào file SQL tương ứng
2. Thêm method vào `CleanupService`
3. Cập nhật `CleanupAllExpiredData` method
4. Thêm config option vào `CronSetting`

### 3. Tắt/bật từng loại cleanup

Trong `config.yaml`:

```yaml
cron:
  enabled: true
  cleanup_sessions: false        # Tắt cleanup sessions
  cleanup_verifications: true   # Bật cleanup verifications
  cleanup_password_resets: true  # Bật cleanup password resets
```

## Best Practices

1. **Monitor Performance**: Theo dõi thời gian thực hiện cleanup
2. **Backup Database**: Đảm bảo có backup trước khi cleanup
3. **Test trong Development**: Test kỹ trước khi deploy production
4. **Log Analysis**: Phân tích logs để tối ưu hóa
5. **Error Alerting**: Thiết lập alert khi có lỗi

## Troubleshooting

### Cron không chạy

1. Kiểm tra `cron.enabled = true` trong config
2. Kiểm tra logs để xem có lỗi khởi tạo không
3. Kiểm tra database connection

### Cleanup thất bại

1. Kiểm tra database permissions
2. Kiểm tra logs để xem lỗi cụ thể
3. Test manual cleanup qua API

### Performance Issues

1. Kiểm tra số lượng records cần cleanup
2. Thêm indexes cho `expires_at` columns
3. Chạy cleanup vào giờ ít traffic

## Security Considerations

1. **Admin Endpoints**: Các endpoint admin cần authentication
2. **Database Permissions**: Chỉ cho phép DELETE trên các bảng cần thiết
3. **Logging**: Không log sensitive data
4. **Rate Limiting**: Áp dụng rate limiting cho admin endpoints

## Kết luận

Cron job cleanup system cung cấp:
- ✅ Tự động dọn dẹp dữ liệu hết hạn
- ✅ Cấu hình linh hoạt
- ✅ Monitoring và logging đầy đủ
- ✅ API quản lý
- ✅ Error handling robust
- ✅ Performance optimized

Hệ thống sẵn sàng cho production với các best practices được áp dụng.

