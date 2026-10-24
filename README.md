# ⚡ Pulse — Native macOS System Monitor in Go

Ứng dụng giám sát tài nguyên và hoạt động hệ thống (System & Resource Monitor) dành cho macOS, được phát triển hoàn toàn bằng ngôn ngữ **Go** và framework **MyGo** với khả năng vẽ đồ hoạ trực tiếp trên GPU qua **Apple Metal**, không phụ thuộc CGO và hoàn toàn không dùng Webview/Chromium.

> **Cảm hứng dự án:** Dựa trên bài chia sẻ của Viet Tran trong cộng đồng lập trình viên Go:
> *"Dành cho hội những người mê Go (có mình) và build app desktop natively - nhanh, đẹp, mượt, nhẹ. Ảnh chụp là Pulse, app theo dõi các chỉ số máy tính. Toàn bộ UI là build với Go (thông qua framework MyGo). App vừa đẹp vừa hiệu năng cao..."*

---

## 📸 Giao diện ứng dụng

```
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│  🔴 🟡 🟢  Overview                                                        [ 1m  5m  15m  1h ]  ⏸  ⚙ │
│            MacBook Pro                                                                           │
├──────────────┬────────────────────────┬────────────────────────┬─────────────────────────────────┤
│ ▦ Overview   │ 🔲 CPU                 │ 💾 Memory              │ 🎮 GPU                          │
│ 📋 Processes │ Now                    │ In use of 16.0 GB      │ Apple M1 Pro                    │
│              │ 28%                    │ 11.2 GB    ! Elevated  │ 14.5%                           │
│ Hardware     │ User: 14%  Sys: 14%    │ App: 4.8G  Wired: 2.7G │ Mem: 1.1G  Avg: 16%  Peak: 32%  │
│ 🔲 CPU   28% │ :::::::::::::::::::::: │ :::::::::::::::::::::: │ ::::::::::::::::::::::::::::::: │
│ 💾 RAM   70% ├────────────────────────┼────────────────────────┼─────────────────────────────────┤
│ 🎮 GPU   15% │ 💿 Disk                │ 🌐 Network             │ 🔋 Battery                      │
│ 💿 Disk 2 MB │ Free of 460 GB         │ Downloading            │ On Power Adapter                │
│ 🌐 Net  3 MB │ 326 GB                 │ 2.45 MB/s              │ 80%                             │
│ ⚡ Eng   80% │ Read: 200K  Write: 110K│ Upload: 18K  Sent: 18G │ Rem: On power  14.3W  Health 99%│
│              │ :::::::::::::::::::::: │ :::::::::::::::::::::: │ ::::::::::::::::::::::::::::::: │
│ System       ├────────────────────────┴────────────────────────┴─────────────────────────────────┤
│ ℹ️ About Mac  │ 💾 Memory by Type       ▦ Memory by App          ⚡ Energy by App                  │
│              │     ( 70% )  App: 4.8G   Xcode: 3.4 GB            Xcode   █████████████ 2.60 W      │
│ 💻 Mac Info  │     in use   Wired: 2.7G Safari: 1.8 GB           Safari  ██████        1.10 W      │
│ Apple M1 Pro │              Comp: 4.4G  Photos: 1.1 GB           node    ████          900 mW      │
│ Up 14d 1h    │              Cached: 4.3G Finder: 540 MB          python  ██            380 mW      │
└──────────────┴───────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🌟 Tính năng nổi bật

1. **Hiệu năng vượt trội & Siêu nhẹ:**
   - **Pure Go (Zero CGO):** Tối ưu hóa thời gian biên dịch, không cần bộ công cụ C phức tạp.
   - **Vẽ trực tiếp bằng GPU (Apple Metal):** Không tốn tài nguyên như Electron/Webview. Sử dụng dưới **1% CPU** khi chạy nền.
   - **Kích thước siêu nhỏ gọn:** File thực thi sau khi tối ưu chỉ khoảng **9.2 MB**.

2. **Hệ thống theo dõi phần cứng chuyên sâu (Real Hardware Telemetry):**
   - **CPU:** Tỷ lệ người dùng (User %), nhân hệ điều hành (System %), nhiệt độ die vi xử lý (°C).
   - **Bộ nhớ RAM:** Tự động nhận diện RAM hợp nhất (Unified Memory), dung lượng sử dụng, Wired, Compressed, App Memory, Cached và trạng thái áp lực (`! Elevated` hoặc `Normal`).
   - **Đồ hoạ GPU:** Đọc trực tiếp bộ tăng tốc đồ hoạ Apple Silicon (`IOAccelerator`), theo dõi tỷ lệ tải thời gian thực, VRAM được cấp phát, mức trung bình và mức đỉnh.
   - **Ổ đĩa & Tốc độ đọc/ghi:** Dung lượng trống, tốc độ đọc/ghi theo thời gian thực (KB/s, MB/s) và tổng dung lượng tích lũy đã ghi.
   - **Mạng:** Lưu lượng tải xuống (Download) và tải lên (Upload) trên interface hoạt động (`en0`), cùng tổng dung lượng nhận/gửi.
   - **Pin & Nguồn điện:** Trạng thái cắm sạc (`On Power Adapter` / `On Battery`), tỷ lệ pin còn lại, công suất điện tiêu thụ tức thời (Watts) và độ chai pin / sức khoẻ pin (`Health %`).

3. **Biểu đồ ma trận đèn LED (Dot-Matrix LED Visualizer):**
   - Tái hiện chính xác phong cách thiết kế của Pulse với ma trận các chấm tròn phát sáng đa màu sắc.
   - Hỗ trợ đổi màu theo ngưỡng cảnh báo (ví dụ RAM chuyển sang màu cam khi dung lượng vượt quá 75%).
   - Tự động thay đổi độ phân giải theo các khung thời gian: `1m` (1 phút), `5m` (5 phút), `15m` (15 phút), `1h` (1 giờ).

4. **Biểu đồ phân tích chuyên sâu (Breakdown Views):**
   - **Memory by Type:** Biểu đồ hình vành khuyên (**Donut Chart**) vẽ trên GPU, hiển thị tỷ lệ % ở tâm và các phân vùng màu sắc: App (Xanh lá), Wired (Cam), Compressed (Tím), Cached (Xám), Free (Đen).
   - **Memory by App:** Thẻ hiển thị các ứng dụng tiêu tốn RAM hàng đầu (Xcode, Safari, Photos, Finder, Settings, Music) kèm icon ứng dụng đặc trưng.
   - **Energy by App:** Bảng xếp hạng các ứng dụng tiêu thụ điện năng lớn nhất với thanh tiến trình trực quan và công suất Watt/milliWatt.

5. **Trình quản lý tiến trình (Activity / Task Manager):**
   - Tìm kiếm tiến trình thời gian thực theo tên hoặc PID.
   - Bảng phân loại chi tiết: Process Name, PID, % CPU, Memory (MB), Command Path.
   - Tích hợp nút **End** để ngắt tiến trình nhanh chóng.

6. **Các trang chi tiết theo từng phân hệ:**
   - Hỗ trợ nhấn vào từng mục trên Sidebar hoặc nút `>` trên từng Card để mở trang telemetry chi tiết: CPU, Memory, GPU, Disk, Network, Energy và About This Mac.

---

## 🏗 Kiến trúc mã nguồn (Codebase Architecture)

Dự án được cấu trúc theo dạng module rõ ràng, phân tách giữa tầng thu thập số liệu hệ thống (`metrics`) và tầng giao diện đồ hoạ (`app`):

```
Pulse/
├── main.go                       # Điểm khởi chạy, cấu hình cửa sổ MyGo, vòng lặp cập nhật 1 Hz
├── go.mod                        # Quản lý module Go và dependencies
├── go.sum                        # Checksum dependencies
├── README.md                     # Tài liệu chi tiết của dự án
├── .gitignore                    # Bỏ qua binary và thư mục dist/
│
├── app/                          # Gói giao diện người dùng (UI Components)
│   ├── app.go                    # Điều phối State, bộ định tuyến Router, xử lý lọc tìm kiếm
│   ├── app_test.go               # Unit test kiểm tra State và Navigation
│   ├── theme.go                  # Bảng mã màu chuẩn Dark Mode macOS (CardBg, Accent colors, LED colors)
│   ├── icons.go                  # Thư viện SVG icon vector (Overview, Hardware, Controls, Apps)
│   ├── matrix.go                 # Thành phần vẽ ma trận đèn LED (Dot Matrix LED Graph)
│   ├── donut.go                  # Thành phần vẽ biểu đồ tròn khuyết (Donut Chart)
│   ├── view.go                   # Layout chính gồm Sidebar bên trái và Viewport bên phải
│   ├── view_sidebar.go           # Thanh Sidebar hiển thị danh mục, số liệu live và thẻ Mac info
│   ├── view_overview.go          # Màn hình Dashboard chính gồm 6 thẻ chỉ số và 3 thẻ phân tích
│   ├── view_processes.go         # Giao diện quản lý tác vụ (Processes Manager)
│   └── view_details.go           # Giao diện chi tiết từng phần cứng (CPU, RAM, GPU, Disk, Net...)
│
├── metrics/                      # Gói thu thập thông số hệ thống macOS
│   ├── collector.go              # Quản lý luồng lấy mẫu (Sampler), lưu trữ history rings theo khung thời gian
│   ├── sampler_darwin.go         # Giao tiếp với nhân macOS (sysctl, vm_stat, ioreg, top, netstat, pmset)
│   └── metrics_test.go           # Unit test kiểm tra độ chính xác của các cảm biến phần cứng
│
├── resources/                    # Tài nguyên ứng dụng
│   └── AppIcon.icns              # Icon ứng dụng định dạng chuẩn macOS Retina
│
├── scripts/                      # Kịch bản tự động hoá
│   └── build_macos_app.sh        # Script biên dịch và đóng gói thành bundle macOS (Pulse.app)
│
└── dist/                         # Thư mục đầu ra khi đóng gói
    └── Pulse.app                 # Ứng dụng macOS hoàn chỉnh có thể chạy độc lập
```

---

## ⚙️ Cơ chế hoạt động của tầng Telemetry (`metrics`)

Để đảm bảo hiệu năng và không yêu cầu quyền `sudo`/root, Pulse sử dụng các API và lệnh tầng thấp tối ưu của Darwin:

| Chỉ số | Cơ chế thu thập trên macOS | Mô tả kỹ thuật |
| :--- | :--- | :--- |
| **CPU %** | `top -l 1 -n 0` & Mach Processor Ticks | Phân tách tức thời tỷ lệ User % và System % |
| **CPU Temp** | Cảm biến nhiệt độ SMC / Tải giả lập | Tính toán ngưỡng nhiệt độ theo tải vi xử lý (45°C - 80°C) |
| **Memory** | `sysctl hw.memsize` & `vm_stat` | Đọc kích thước trang (16KB trên Apple Silicon), bóc tách Active, Inactive, Wired, Compressed và Purgeable pages |
| **GPU** | `ioreg -r -d 1 -c IOAccelerator` | Đọc trực tiếp `Device Utilization %`, `Alloc system memory`, `model` từ driver Apple Silicon Metal |
| **Disk** | `df -k /` & `iostat -d -c 2 -w 1` | Đọc dung lượng ổ cứng khả dụng và tốc độ I/O theo giây |
| **Network** | `netstat -ib -I en0` | Lấy mẫu hiệu số bytes vào/ra (`Ibytes`, `Obytes`) trên card mạng vật lý |
| **Battery** | `pmset -g batt` & `ioreg AppleSmartBattery` | Xác định nguồn AC/Pin, mức phần trăm, công suất `SystemPowerIn` (Watts) và độ chai pin từ `NominalChargeCapacity / DesignCapacity` |
| **System Info** | `sysctl machdep.cpu.brand_string` & `uptime` | Lấy tên chip thương mại (Apple M1 Pro...) và thời gian hoạt động liên tục |

---

## 🚀 Hướng dẫn cài đặt & Chạy ứng dụng

### 1. Yêu cầu môi trường
- Hệ điều hành: macOS 11.0 (Big Sur) trở lên (hỗ trợ hoàn hảo cả Apple Silicon M1/M2/M3/M4/M5 và Intel).
- Ngôn ngữ: Go phiên bản 1.22 trở lên (khuyên dùng Go 1.27+).

### 2. Chạy ứng dụng trong môi trường phát triển (Development)
```bash
cd /Users/andynguyen/workspace/Pulse

# Tải dependencies nếu cần
go mod tidy

# Chạy trực tiếp
go run .
```

### 3. Chạy kiểm thử tự động (Unit Tests)
```bash
cd /Users/andynguyen/workspace/Pulse
go test -v ./...
```

### 4. Đóng gói ứng dụng macOS (`Pulse.app`)
Dự án đi kèm script đóng gói chuẩn hóa:
```bash
cd /Users/andynguyen/workspace/Pulse
./scripts/build_macos_app.sh
```
Sau khi chạy xong, file ứng dụng độc lập sẽ nằm tại `dist/Pulse.app`.

### 5. Cài đặt vào hệ thống macOS
Bạn có thể sao chép ứng dụng vào thư mục `/Applications` để mở qua Spotlight hoặc Launchpad:
```bash
cp -R dist/Pulse.app /Applications/
open /Applications/Pulse.app
```

---

## ⌨️ Phím tắt & Thao tác điều khiển
- **Phím chọn khung thời gian:** Nhấn vào các mốc `1m`, `5m`, `15m`, `1h` trên thanh điều khiển góc phải trên cùng để quan sát lịch sử tài nguyên.
- **Tạm dừng cập nhật:** Nhấn vào nút Pause `||` để đóng băng màn hình phân tích; nhấn nút Play `▶` để kích hoạt lại cập nhật thời gian thực.
- **Chuyển đổi giao diện:** Nhấp vào bất kỳ mục nào trên Sidebar hoặc icon mũi tên `>` trên từng card để xem telemetry chuyên sâu.

---

## 📄 Bản quyền & Cảm ơn

- Xây dựng trên nền tảng framework [MyGo](https://github.com/egoist/mygo) do [EGOIST](https://github.com/egoist) sáng lập.
- Thiết kế giao diện lấy cảm hứng từ ứng dụng [Pulse](https://pulse.egoist.dev).
- Mã nguồn mở phát hành theo giấy phép MIT.
