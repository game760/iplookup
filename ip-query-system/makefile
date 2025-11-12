# 基础配置
BINARY_LINUX = ip-query-service  # Linux 可执行文件名（无后缀）
BINARY_WIN = ip-query-service.exe  # Windows 可执行文件名（带 .exe 后缀）
GOARCH ?= amd64  # 架构（默认 64位）

# 默认目标：同时编译 Linux 和 Windows 版本
all: build-linux build-windows

# 编译 Linux 版本（默认架构 amd64）
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -o $(BINARY_LINUX) main.go
	@echo "Linux 版本编译完成: $(BINARY_LINUX)"

# 编译 Windows 版本（默认架构 amd64）
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=$(GOARCH) go build -o $(BINARY_WIN) main.go
	@echo "Windows 版本编译完成: $(BINARY_WIN)"

# 单独编译 Linux 版本（快捷命令）
linux: build-linux

# 单独编译 Windows 版本（快捷命令）
win: build-windows

# 清理所有编译产物
clean:
	rm -f $(BINARY_LINUX) $(BINARY_WIN)
	@echo "已清理所有编译文件"

# 运行当前系统的可执行文件（自动判断系统）
run:
	@if [ "$(shell uname)" = "Linux" ]; then \
		./$(BINARY_LINUX); \
	elif [ "$(shell uname -s | cut -d' ' -f1)" = "MINGW64" ] || [ "$(shell uname -s)" = "Windows_NT" ]; then \
		./$(BINARY_WIN); \
	else \
		echo "不支持的系统"; \
	fi