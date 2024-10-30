BTL 1:
## Hướng dẫn cài đặt

Lưu ý đổi tên file từ a.png sang file để download. Đã test với file png, txt chưa thử với các file khác. 

File split.go ở folder src dùng để tách file thành các pieces. ``` go run split.go```

File client.go ở folder client dùng để khi download thì sẽ tải từ các client này. Chạy lệnh ``` go run client.go 5000``` và ``` go run client.go 5001```

File main.go ở folder download_client dùng để tải pieces và merge thành file. Lệnh chạy ``` go run main.go```

