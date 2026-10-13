#!/bin/sh
set -eu

GOBIN=${BUILD_TOOLS}
BIN="${GOBIN}/mockgen"

if [ -f "${BIN}" ]; then
  INSTALLED="$("${BIN}" -version 2>/dev/null || true)"
  EXPECTED="${VERSION}"
  if [ "x${INSTALLED}" == "x${EXPECTED}" ]; then
    echo "mockgen ${VERSION} 已安装"
    exit 0
  fi

  echo "mockgen 版本不匹配（现存：${INSTALLED:-无}, 需要：${EXPECTED}），重新安装..."
  mv "${BIN}" "${BIN}.bk"
fi

echo "正在安装 mockgen..."
GOBIN=${GOBIN} go install "go.uber.org/mock/mockgen@${VERSION}"

[ -f "${BIN}.bk" ] && rm "${BIN}.bk" || true
