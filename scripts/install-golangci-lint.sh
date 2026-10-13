#!/bin/sh
set -eu

# 提取版本号中的数字部分（去掉可能的 v 前缀）
# 输入可能形如 "v0.4.0"、"0.4.0"、"mockgen v0.4.0" 等
# 只取第一个 x.y.z
normalize_version() {
  if [ "$#" -gt 0 ]; then
    printf '%s' "$1"
  elif [ ! -t 0 ]; then
    cat
  fi |
    grep -oE 'v?[0-9]+\.[0-9]+\.[0-9]+' |
    head -n 1 |
    sed 's/^v//'
}

if [ -z "$BUILD_TOOLS" ] || [ -z "$VERSION" ]; then
  echo "Both BUILD_TOOLS and VERSION environment required."
  exit 1
fi

BIN="${BUILD_TOOLS}/golangci-lint"

if [ -f "${BIN}" ]; then
  INSTALLED="$("${BIN}" version 2>/dev/null | normalize_version || true)"
  EXPECTED=$(normalize_version "${VERSION}")
  if [ "x${INSTALLED}" == "x${EXPECTED}" ]; then
    echo "golangci-lint ${VERSION} 已安装"
    exit 0
  fi

  echo "golangci-lint 版本不匹配（现存：${INSTALLED:-无}, 需要：${EXPECTED}），重新安装..."
  mv "${BIN}" "${BIN}.bk"
fi

echo "正在安装 golangci-lint..."
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ${BUILD_TOOLS} v${VERSION}

[ -f "${BIN}.bk" ] && rm "${BIN}.bk" || true
