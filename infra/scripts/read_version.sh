#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/../.." && pwd)"
version_file="${repo_root}/VERSION"

if [[ ! -f "${version_file}" ]]; then
  echo "VERSION file not found: ${version_file}" >&2
  exit 1
fi

version="$(tr -d '[:space:]' < "${version_file}")"

if [[ -z "${version}" ]]; then
  echo "VERSION file is empty: ${version_file}" >&2
  exit 1
fi

printf '{"version":"%s"}\n' "${version}"
