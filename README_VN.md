<h1 align="center">Ichain </h1>
<h4 align="center">Version 0.0.1 </h4>

Vietnameses | [English](README.md)

Chào mứng bạn đến với thư viên mã nguồn Ichain

Ichain sử dụng blockchain trong truy xuất chuỗi cung ứng. Ichain giúp triển khai và chạy ứng dụng truy xuất chuỗi cung ứng trên blockchain dễ dàng hơn.

Mã nguồn hiện đang ở giai đoạn thử nghiệm, đang trong quá trình phát triển nhanh chóng. Mã chính hiện tại không ổn định, các phiên bản ổn định sẽ được liệt kê trên trang phát hành.

Bất kỳ câu hỏi liên quan đến việc hợp tác triển khai ứng dụng xin vui lòng gửi đến email (hotro@icheck.vn)


#### Các tính năng
- Tốc độ xử  lý giao dịch ngay lập tức.
- Thời gian tạo khối nhanh.
- Truy xuất nguồn gốc chuỗi cung ứng.
- Triển khai và quản lý dòng chảy của sản phẩm.
- Có thể mở rộng hợp đồng thông minh.
- Trao đổi tài sản nhanh chóng.

#### Bắt đầu

```
$ make
````

Sau khi xây dựng mã nguồn thành công, bạn sẽ thấy 2 chương trình thực thị trong thư mục ./build

- `ichaind`: Chương trình dòng lệnh để triển khai nút
- `ichaincli`: Chương trình dòng lệnh khách hàng thực thi giao dịch 


# Triển khai nút thử nghiệm nút công cộng

1. Tạo tài khoản
- Thông qua chương trình dòng lệnh, khởi tạo một tài khoản
```
./ichaincli keys add testaccount
Enter a passphrase for your key:
Re-enter password:
testaccount     283873F09FEBC7EC95BCFBD43B37CF0678B8232A
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

melody tunnel slice calm basket round retreat cry impulse tail tunnel awkward morning wash apple abandon
```
2. Bắt đầu ichain
- Thông qua chương trình dòng lệnh, bắt đầu nút công cộng:
```
./ichaind start
```

### Thực thi

Chạy `ichaincli --help` để  để xem hướng dẫn chi tiết

### Ví dụ Trao đổi tài sản
```
./ichaincli transfer --name testaccount --amount 100tomato --to 283873F09FEBC7EC95BCFBD43B37CF0678B8232A
```

### Đóng góp
Mọi thông tin đóng góp về dự án xin vui lòng gửi email đến đia chỉ (hotro@icheck.vn)