<h1>Ichain </h1>
<h4>Version 0.0.1 </h4>

Branch    | Tests | Coverage
----------|-------|---------
develop   | [![CircleCI](https://circleci.com/gh/icheckteam/ichain/tree/develop.svg?style=shield)](https://circleci.com/gh/icheckteam/ichain/tree/develop) | [![codecov](https://codecov.io/gh/icheckteam/ichain/branch/develop/graph/badge.svg)](https://codecov.io/gh/icheckteam/ichain)
master    | [![CircleCI](https://circleci.com/gh/icheckteam/ichain/tree/master.svg?style=shield)](https://circleci.com/gh/icheckteam/ichain/tree/master) | [![codecov](https://codecov.io/gh/icheckteam/ichain/branch/master/graph/badge.svg)](https://codecov.io/gh/icheckteam/ichain)

Vietnameses | [English](README.md)

Chào mừng bạn đến với thư viện mã nguồn Ichain

Ichain là một blockchain được phát triển dựa trên tendermint giúp triển khai, k

Mã nguồn hiện đang ở giai đoạn thử nghiệm, đang trong quá trình phát triển nhanh chóng. Mã chính hiện tại không ổn định, các phiên bản ổn định sẽ được liệt kê trên trang phát hành.

Bất kỳ câu hỏi liên quan đến việc hợp tác triển khai ứng dụng xin vui lòng gửi đến email (dev@icheck.vn)


#### Các tính năng
- Tốc độ xử  lý giao dịch ngay lập tức.
- Thời gian tạo khối nhanh.
- Truy xuất nguồn gốc chuỗi cung ứng.
- Triển khai và quản lý dòng chảy của sản phẩm.
- Có thể mở rộng hợp đồng thông minh.
- Trao đổi tài sản nhanh chóng.
- Nhận dạng danh tính kỹ thuật số.

## Các module

1. [Identity](https://github.com/icheckteam/documentation/blob/master/Identity.md) là một module quản lý và nhận dạng danh tính kỹ thuật số.
2. [Asset](https://github.com/icheckteam/documentation/blob/master/Asset.md) là một module quản lý, và trao đổi tài sản ký thuật số.

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

### Trao đổi tài sản
```
./ichaincli transfer --name testaccount --amount 100tomato --to 283873F09FEBC7EC95BCFBD43B37CF0678B8232A
```
### Đóng góp
Mọi thông tin đóng góp về dự án xin vui lòng gửi email đến đia chỉ (dev@icheck.vn)
