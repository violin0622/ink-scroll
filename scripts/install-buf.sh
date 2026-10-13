#!/bin/sh
set -eu

if [ -z "$BUILD_TOOLS" ] || [ -z "$VERSION" ]; then
  echo "Both BUILD_TOOLS and VERSION environment required."
  exit 1
fi

BIN="${BUILD_TOOLS}/buf"

if [ -f "${BIN}" ]; then
  INSTALLED="$("${BIN}" --version 2>/dev/null || true)"
  EXPECTED="${VERSION}"
  if [ "x${INSTALLED}" == "x${EXPECTED}" ]; then
    echo "buf ${VERSION} 已安装"
    exit 0
  fi

  echo "buf 版本不匹配（现存：${INSTALLED:-无}, 需要：${EXPECTED}），重新安装..."
  mv "${BIN}" "${BIN}.bk"
fi

echo "正在安装 buf..."
GOBIN=${BUILD_TOOLS} go install github.com/bufbuild/buf/cmd/buf@v${VERSION}

[ -f "${BIN}.bk" ] && rm "${BIN}.bk" || true
