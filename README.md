# Deploy Go - Dự án Golang cơ bản

Đây là một dự án Golang cơ bản sử dụng Gin framework, được thiết kế để deploy lên Render.

## Tính năng

- Web server đơn giản với Gin framework
- API endpoints cơ bản:
  - `GET /` - Trang chủ
  - `GET /health` - Health check
  - `GET /api/hello` - API endpoint mẫu

## Cài đặt và chạy local

1. Cài đặt Go (phiên bản 1.21 trở lên)
2. Clone repository này
3. Chạy lệnh:
   ```bash
   go mod download
   go run main.go
   ```
4. Truy cập http://localhost:8080

## Deploy lên Render

1. Push code lên GitHub repository
2. Đăng nhập vào Render.com
3. Tạo Web Service mới
4. Connect với GitHub repository
5. Cấu hình:
   - Build Command: `go mod download && go build -o main .`
   - Start Command: `./main`
   - Environment: `Go`
   - Health Check Path: `/health`

Hoặc sử dụng file `render.yaml` đã được cấu hình sẵn.

## API Endpoints

- `GET /` - Trang chủ với thông báo chào mừng
- `GET /health` - Health check endpoint
- `GET /api/hello` - API endpoint mẫu

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - Web framework cho Go
