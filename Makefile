VERSION = v0.1.0
PHASE = beta
FULL_VERSION := $(shell \
	BASE_VERSION=$$(git describe --exact-match --tags HEAD 2>/dev/null || echo $(VERSION)); \
	if [ -n "$(PHASE)" ]; then \
		echo "$$BASE_VERSION-$(PHASE)+$$(git rev-parse --short HEAD)"; \
	else \
		echo "$$BASE_VERSION"; \
	fi) 

BUF_VERSION = 1.59.0
GOLANGCI_LINT_VERSION = 2.12.2
MIGRATE_VERSION = v4.15.2
MOCKGEN_VERSION = v0.6.0

# 构建信息变量
GIT_SHA := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

# 探测本地操作系统和架构
OS := $(shell uname -s)
ARCH := $(shell uname -m)

# 数据库迁移变量
MIGRATION ?=

.PHONY: build
build: gen
	@echo "Building version: $(FULL_VERSION)"; 
	go build -ldflags "\
		-X play-ddd/cmd.Version=$(FULL_VERSION) \
		-X play-ddd/cmd.GitSHA=$(GIT_SHA) \
		-X play-ddd/cmd.BuiltAt=$(BUILD_TIME) \
		-X play-ddd/cmd.GoVersion=$(GO_VERSION)" \
		-o build/app .

.PHONY: run
run: build
	@./bin/app

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/app
	@go clean -testcache

BUILD_TOOLS = $(CURDIR)/build/tools
$(BUILD_TOOLS):
	@mkdir -p $(BUILD_TOOLS)

.PHONY: golangci-lint
golangci-lint: $(BUILD_TOOLS)
	@BUILD_TOOLS=$(BUILD_TOOLS) VERSION=$(GOLANGCI_LINT_VERSION) \
		sh scripts/install-golangci-lint.sh

.PHONY: buf
buf: $(BUILD_TOOLS)
	@BUILD_TOOLS=$(BUILD_TOOLS) VERSION=$(BUF_VERSION) \
		sh scripts/install-buf.sh

.PHONY: migrate
migrate: $(BUILD_TOOLS)
	@MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(OS)-$(ARCH).tar.gz"; \
	MIGRATE_BIN="$(BUILD_TOOLS)/migrate"; \
	if [ ! -f $$MIGRATE_BIN ]; then \
		echo "未找到 migrate，正在下载..."; \
		curl -sSL "$$MIGRATE_URL" | tar -xz -C $(BUILD_TOOLS) migrate; \
	else \
		INSTALLED=$$($$MIGRATE_BIN -version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo ""); \
		EXPECTED=$$(echo $(MIGRATE_VERSION) | sed 's/v//'); \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "migrate 版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新下载..."; \
			curl -sSL "$$MIGRATE_URL" | tar -xz -C $(BUILD_TOOLS) migrate; \
		fi; \
	fi

.PHONY: mockgen
mockgen: $(BUILD_TOOLS)
	@BUILD_TOOLS=$(BUILD_TOOLS) VERSION=$(MOCKGEN_VERSION) \
	sh scripts/install-mockgen.sh

.PHONY: fmt
fmt: golangci-lint buf
	@$(BUILD_TOOLS)/golangci-lint fmt
	@$(BUILD_TOOLS)/buf format -w

.PHONY: lint
lint: fmt golangci-lint
	@$(BUILD_TOOLS)/golangci-lint run --fix
	@$(BUILD_TOOLS)/buf lint

.PHONY: gen
gen: buf
	@$(BUILD_TOOLS)/buf generate

.PHONY: schema-new
schema-new: migrate
	@if [ -z "$(MIGRATION)" ]; then \
		echo "错误: 请提供迁移名称，例如: make schema-new MIGRATION=my_migration" && exit 1; \
	fi; 
	@$(CURDIR)/bin/migrate create -ext sql -dir $(CURDIR)/schemas/migrations -seq $(MIGRATION)

.PHONY: schema-up
schema-up: n:=
schema-up: migrate
	$(CURDIR)/bin/migrate -database ${DB} -path $(CURDIR)/schemas/migrations up $(n)

.PHONY: schema-down
schema-down: migrate
	@$(CURDIR)/bin/migrate -database ${DB} -path $(CURDIR)/schemas/migrations down
