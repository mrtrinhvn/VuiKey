# 📱 Hồ Sơ Đưa VuiKey Lên Google Play Store (CH Play)

Tài liệu này chuẩn bị sẵn toàn bộ nội dung метаdata, bảng khai báo **An toàn dữ liệu (Data Safety)** và hướng dẫn xuất bản file **`.aab` (`VuiKey-1.0.0-android-playstore.aab`)** lên **Google Play Console**.

---

## 1. Thông Tin Kỹ Thuật Chuẩn Google Play Console (API 35)
- **Application ID (Package Name):** `com.echipco.vuikey`
- **Target SDK:** `35` (Android 15 — Đạt chuẩn bắt buộc mới nhất của Google Play Console)
- **Min SDK:** `24` (Hỗ trợ từ Android 7.0 trở lên — bao phủ >99.5% thiết bị Android đang hoạt động)
- **Quyền truy cập (`Permissions`):** Chỉ duy nhất `android.permission.VIBRATE` (rung phản hồi phím). **Hoàn toàn KHÔNG yêu cầu quyền `INTERNET`** — giúp VuiKey vượt qua quy trình xét duyệt bảo mật bàn phím của Google Play cực nhanh và đạt nhãn *"Không thu thập bất kỳ dữ liệu nào"*.
- **File tải lên Google Play Console:** `dist/VuiKey-1.0.0-android-playstore.aab`
- **File cài đặt trực tiếp (APK):** `dist/VuiKey-1.0.0-android.apk`

---

## 2. Nội Dung Trang Cửa Hàng (Store Listing — Tiếng Việt `vi-VN`)

### Tên ứng dụng (App Name — Tối đa 30 ký tự):
```text
VuiKey: Bộ Gõ Tiếng Việt 0.0ms
```

### Mô tả ngắn (Short Description — Tối đa 80 ký tự):
```text
Bàn phím Tiếng Việt không gạch chân, tự hiểu Tiếng Anh & sửa dấu Backspace ngay!
```

### Mô tả đầy đủ (Full Description — Tối đa 4000 ký tự):
```text
Bộ gõ tiếng Việt ngày nay không còn là thứ gì đao to búa lớn — nhưng VuiKey lập tức tạo ra khác biệt khi giải quyết dứt điểm những "nỗi đau" tồn tại suốt nhiều năm trên điện thoại và máy tính bảng mà chưa bàn phím nào xử lý trọn vẹn.

❌ 4 NỖI ĐAU KINH ĐIỂN KHI DÙNG CÁC BÀN PHÍM KHÁC:
• Đang gõ Tiếng Việt mà gõ Tiếng Anh là hỏng chữ: Gõ window, view, test, user, server, linux... là bị biến thành ưindoư, vieư, tét, úer, lĩnu — buộc phải xóa đi và bấm đổi ngôn ngữ liên tục.
• Gạch chân khó chịu & nhảy đúp chữ: Bộ đệm gạch chân (preedit) trên các ứng dụng nhắn tin và thanh địa chỉ trình duyệt Chrome thường xuyên gây nhân đôi cụm chữ hoặc nuốt chữ cái đầu.
• Sửa 1 ký tự phải xóa sạch cả từ: Lỡ chạm nhầm 1 phím rồi bấm Backspace xóa đi, nhưng gõ tiếp dấu sắc/huyền/mũ/móc thì bàn phím không nhận nữa, bắt phải xóa lùi về tận chữ đầu tiên.
• Gõ nhanh bị rớt dấu: Bắt phải gõ dd ngay đầu từ hoặc ép đúng 1 thứ tự gõ ươ.

⚡ VUIKEY GIẢI QUYẾT THẾ NÀO? — NHANH TỨC THÌ & THÔNG MINH THỰC SỰ:
1. Không gạch chân & Chống nhân đôi chữ 100% (Zero-Preedit Batch Edit): Chèn chữ nguyên tử với độ trễ 0.00ms, không bao giờ hiện gạch chân dưới từ đang gõ, tương thích tuyệt đối với mọi khung chat (Zalo, Messenger, Telegram, Chrome, Docs).
2. Tự hiểu Tiếng Anh mà không cần tắt Tiếng Việt (Smart English Guard): Phân tích cấu trúc âm vị học thời gian thực để tự động giữ nguyên 100% các từ tiếng Anh (window, test, function, server, linux...) ngay khi đang bật chế độ Tiếng Việt.
3. Bấm Backspace 1 nhịp sửa dấu ăn ngay (Instant Unseal): Xóa lùi 1 ký tự là mở khóa âm tiết ngay lập tức để bạn bỏ dấu hoặc sửa vần tiếp, không bao giờ phải xóa lại từ đầu.
4. Gõ tự do mọi thứ tự & Tự tách âm tiết viết liền (Free-Style Telex/VNI + Sliding Syllable Window):
   - Cho phép gõ chữ d thứ hai ở bất kỳ đâu trong từ (dodongj, doongjd ➔ động).
   - Nhận trọn cả 5 cách gõ nguyên âm đôi ươ (dduowcj, dduocwj, dduwowcj, dduwocwj, dduwocj ➔ được).
   - Tự động tách từng âm tiết kể cả khi viết dính liền không dấu cách (nhacjdiujdangf ➔ nhạcdịudàng).
5. Vuốt phím Spacebar làm Trackpad di chuyển con trỏ mượt mà + Hỗ trợ bàn phím cứng Bluetooth / Samsung DeX (chuyển nhanh Việt/Anh bằng Ctrl + Space).
6. Bảo mật tuyệt đối 100% Offline: VuiKey KHÔNG yêu cầu quyền truy cập Internet (0 Network Permission). Mọi thao tác gõ phím của bạn đều được xử lý hoàn toàn trên thiết bị.

Phát triển bởi Công ty TNHH eChip (echipco.com) & VuiCaHat (vuicahat.com).
```

---

## 3. Khai Báo An Toàn Dữ Liệu (Data Safety) & Chính Sách Bảo Mật
- **URL Chính sách bảo mật (Privacy Policy URL):** `https://echipco.com/chinh-sach.html#bao-mat`
- **Ứng dụng có thu thập hoặc chia sẻ bất kỳ loại dữ liệu người dùng bắt buộc nào không?** ➔ Chọn **KHÔNG (No)** (vì ứng dụng không hề có quyền `android.permission.INTERNET` trong `AndroidManifest.xml`).
- **Quảng cáo (Ads):** Chọn **Không, ứng dụng của tôi không chứa quảng cáo**.
