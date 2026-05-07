#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPT_DIR}/build"
PACKAGE_DIR="${BUILD_DIR}/package"
PYTHON_VERSION="${PYTHON_VERSION:-3.12}"
PLATFORM="${PLATFORM:-manylinux2014_x86_64}"

rm -rf "${PACKAGE_DIR}"
mkdir -p "${PACKAGE_DIR}"

python3 -m pip install \
  --platform "${PLATFORM}" \
  --implementation cp \
  --python-version "${PYTHON_VERSION}" \
  --only-binary=:all: \
  --upgrade \
  --target "${PACKAGE_DIR}" \
  -r "${SCRIPT_DIR}/requirements.txt"

cp "${SCRIPT_DIR}/app.py" "${PACKAGE_DIR}/app.py"

find "${PACKAGE_DIR}" -type d -name "__pycache__" -prune -exec rm -rf {} +
find "${PACKAGE_DIR}" -type f \( -name "*.pyc" -o -name "*.pyo" \) -delete

echo "Built Lambda package directory: ${PACKAGE_DIR}"
