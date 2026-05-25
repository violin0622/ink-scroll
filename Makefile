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
GOLANGCI_LINT_VERSION = v2.12.2
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
		-o bin/app .

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

$(CURDIR)/bin:
	@mkdir -p $(CURDIR)/bin

.PHONY: golangci-lint
golangci-lint: $(CURDIR)/bin
	@CMD="curl -sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
		| sh -s -- -b $(CURDIR)/bin $(GOLANGCI_LINT_VERSION)"; \
	if [ ! -f $(CURDIR)/bin/golangci-lint ]; then \
		echo "未找到 golangci-lint，正在安装..."; \
		eval $$CMD; \
	else \
		INSTALLED=$$($(CURDIR)/bin/golangci-lint version 2>/dev/null \
			| grep -oE 'version [0-9.]+' | grep -oE '[0-9.]+' || echo ""); \
		EXPECTED=$$(echo $(GOLANGCI_LINT_VERSION) | sed 's/v//'); \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "golangci-lint 版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新安装..."; \
			eval $$CMD; \
		else \
			echo "golangci-lint $(GOLANGCI_LINT_VERSION) 已安装"; \
		fi; \
	fi 

.PHONY: buf
buf: $(CURDIR)/bin
	@BUF_URL="https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$(OS)-$(ARCH)"; \
	BUF_BIN="$(CURDIR)/bin/buf"; \
	if [ ! -f $$BUF_BIN ]; then \
		echo "未找到 buf，正在安装..."; \
		curl -sSL "$$BUF_URL" -o $$BUF_BIN; \
		chmod +x $$BUF_BIN; \
	else \
		INSTALLED=$$($$BUF_BIN --version 2>/dev/null \
			| grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo ""); \
		EXPECTED="$(BUF_VERSION)"; \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "buf 版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新安装..."; \
			curl -sSL "$$BUF_URL" -o $$BUF_BIN; \
			chmod +x $$BUF_BIN; \
		else \
			echo "buf $(BUF_VERSION) 已安装"; \
		fi; \
	fi

.PHONY: migrate
migrate: $(CURDIR)/bin
	@MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(OS)-$(ARCH).tar.gz"; \
	MIGRATE_BIN="$(CURDIR)/bin/migrate"; \
	if [ ! -f $$MIGRATE_BIN ]; then \
		echo "未找到 migrate，正在下载..."; \
		curl -sSL "$$MIGRATE_URL" | tar -xz -C $(CURDIR)/bin migrate; \
	else \
		INSTALLED=$$($$MIGRATE_BIN -version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1 || echo ""); \
		EXPECTED=$$(echo $(MIGRATE_VERSION) | sed 's/v//'); \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "migrate 版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新下载..."; \
			curl -sSL "$$MIGRATE_URL" | tar -xz -C $(CURDIR)/bin migrate; \
		fi; \
	fi

.PHONY: mockgen
mockgen: $(CURDIR)/bin
	@MOCKGEN_BIN="$(CURDIR)/bin/mockgen"; \
	if [ ! -f $$MOCKGEN_BIN ]; then \
		echo "未找到 mockgen，正在安装..."; \
		GOBIN=$(CURDIR)/bin go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION); \
	else \
		INSTALLED=$$($$MOCKGEN_BIN -version 2>/dev/null \
			| grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' | head -1 | sed 's/^v//' || echo ""); \
		EXPECTED=$$(echo $(MOCKGEN_VERSION) | sed 's/v//'); \
		if [ "$$INSTALLED" != "$$EXPECTED" ] || [ -z "$$INSTALLED" ]; then \
			echo "mockgen 版本不匹配 (已安装: $$INSTALLED, 需要: $$EXPECTED)"; \
			echo "正在重新安装..."; \
			GOBIN=$(CURDIR)/bin go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION); \
		else \
			echo "mockgen $(MOCKGEN_VERSION) 已安装"; \
		fi; \
	fi

.PHONY: fmt
fmt: golangci-lint buf
	@$(CURDIR)/bin/golangci-lint fmt
	@$(CURDIR)/bin/buf format -w

.PHONY: lint
lint: fmt golangci-lint
	@$(CURDIR)/bin/golangci-lint run --fix
	@$(CURDIR)/bin/buf lint

.PHONY: gen
gen: buf
	@$(CURDIR)/bin/buf generate

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
