.PHONY: build clean install test

# 默认目标
all: build

# 构建
build:
	@echo "Building skillkit..."
	@go mod tidy
	@go build -o bin/skillkit ./cmd
	@ln -sf skillkit bin/sk
	@chmod +x skillkit sk
	@echo "Done. Run './skillkit' or './sk' to start."

# 清理
clean:
	@rm -rf bin/skillkit bin/sk
	@echo "Cleaned."

# 安装到系统
install: build
	@echo "Installing to /usr/local/bin..."
	@sudo cp bin/skillkit /usr/local/bin/skillkit
	@sudo ln -sf skillkit /usr/local/bin/sk
	@echo "Installed. Run 'skillkit' or 'sk' from anywhere."

# 卸载
uninstall:
	@sudo rm -f /usr/local/bin/skillkit /usr/local/bin/sk
	@echo "Uninstalled."

# 测试
test:
	@go test -v ./...

# 初始化用户配置目录
init:
	@mkdir -p ~/.config/agent/skill
	@mkdir -p ~/.config/agent/agent
	@cp platforms.toml ~/.config/agent/
	@echo "Initialized ~/.config/agent/"

# 帮助
help:
	@echo "Available targets:"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install to /usr/local/bin"
	@echo "  make uninstall- Remove from /usr/local/bin"
	@echo "  make init     - Initialize ~/.config/agent/"
	@echo "  make test     - Run tests"
