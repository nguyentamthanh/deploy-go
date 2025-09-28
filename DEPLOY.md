# Hướng dẫn Deploy lên Render

## Bước 1: Chuẩn bị Repository

1. Tạo repository mới trên GitHub
2. Push code lên repository:
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin https://github.com/yourusername/your-repo-name.git
   git push -u origin main
   ```

## Bước 2: Deploy trên Render

### Cách 1: Sử dụng Render Dashboard

1. Đăng nhập vào [Render.com](https://render.com)
2. Click "New +" → "Web Service"
3. Connect với GitHub và chọn repository
4. Cấu hình:
   - **Name**: `deploy-go` (hoặc tên bạn muốn)
   - **Environment**: `Go`
   - **Build Command**: `go mod download && go build -o main .`
   - **Start Command**: `./main`
   - **Plan**: Free (hoặc Starter/Pro nếu bạn có nhu cầu)
5. Click "Create Web Service"

### Cách 2: Sử dụng render.yaml (Auto-deploy)

1. File `render.yaml` đã được cấu hình sẵn
2. Khi push code lên GitHub, Render sẽ tự động deploy
3. Render sẽ đọc file `render.yaml` và cấu hình service

## Bước 3: Cấu hình Environment Variables

Trong Render Dashboard:

- **GIN_MODE**: `release` (để chạy ở production mode)

## Bước 4: Health Check

Render sẽ tự động kiểm tra health endpoint tại `/health`

## Bước 5: Custom Domain (Optional)

1. Trong Render Dashboard, vào Settings của service
2. Add Custom Domain
3. Cấu hình DNS theo hướng dẫn của Render

## Troubleshooting

### Lỗi Build

- Kiểm tra Go version (cần 1.21+)
- Đảm bảo `go.mod` và `go.sum` đã được commit

### Lỗi Runtime

- Kiểm tra logs trong Render Dashboard
- Đảm bảo PORT environment variable được set (Render tự động set)

### Lỗi Health Check

- Kiểm tra endpoint `/health` hoạt động
- Đảm bảo server start thành công

## Monitoring

- Sử dụng Render Dashboard để monitor logs
- Có thể setup alerts cho downtime
- Free plan có giới hạn 750 giờ/tháng

## Scaling

- Free plan: 1 instance
- Starter plan: Auto-scaling
- Pro plan: Multiple regions

## Cost

- **Free**: $0/tháng (750 giờ, sleep sau 15 phút không hoạt động)
- **Starter**: $7/tháng (không sleep)
- **Pro**: $25/tháng (multiple regions)
