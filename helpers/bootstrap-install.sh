#!/usr/bin/env bash
set -euo pipefail

TRUSTSSH_REPO="${TRUSTSSH_REPO:-https://github.com/nikon-63/TrustSSH.git}"
TRUSTSSH_BRANCH="${TRUSTSSH_BRANCH:-main}"
TRUSTSSH_ENDPOINT="${TRUSTSSH_ENDPOINT:-}"

WORKDIR="/tmp/trustssh-bootstrap"
BOOTSTRAP_DIR="$WORKDIR/TrustSSH/bootstrap-tooling"

log() {
    echo "[trustssh] $*"
}

fail() {
    echo "[trustssh] ERROR: $*" >&2
    exit 1
}

require_root() {
    if [ "$(id -u)" -ne 0 ]; then
        fail "This installer must be run as root."
    fi
}

detect_package_manager() {
    if command -v apt-get >/dev/null 2>&1; then
        PKG_MANAGER="apt"
    elif command -v apk >/dev/null 2>&1; then
        PKG_MANAGER="apk"
    elif command -v dnf >/dev/null 2>&1; then
        PKG_MANAGER="dnf"
    elif command -v yum >/dev/null 2>&1; then
        PKG_MANAGER="yum"
    else
        fail "Unsupported package manager."
    fi
}

install_packages() {
    log "Installing prerequisites"

    case "$PKG_MANAGER" in
        apt)
        apt-get update
        DEBIAN_FRONTEND=noninteractive apt-get install -y python3 ansible git openssh-server
        ;;
        apk)
        apk update
        apk add --no-cache python3 ansible git openssh
        ;;
        dnf)
        dnf install -y python3 ansible git openssh-server
        ;;
        yum)
        yum install -y python3 ansible git openssh-server
        ;;
    esac
}

validate_endpoint() {
    if [ -z "$TRUSTSSH_ENDPOINT" ]; then
        fail "TRUSTSSH_ENDPOINT is not set."
    fi

    TRUSTSSH_ENDPOINT="${TRUSTSSH_ENDPOINT%/}"
}

fetch_public_key() {
    log "Fetching TrustSSH CA public key from $TRUSTSSH_ENDPOINT/public_key.txt"

    TRUSTSSH_CA_PUBLIC_KEY="$(curl --proto '=https' --tlsv1.2 -sSfL "$TRUSTSSH_ENDPOINT/public_key.txt" | tr -d '\r' | sed -e 's/[[:space:]]*$//')"

    if ! printf '%s\n' "$TRUSTSSH_CA_PUBLIC_KEY" | grep -Eq '^ssh-(ed25519|rsa) [A-Za-z0-9+/=]+( .*)?$'; then
        fail "Downloaded public key is not a valid SSH public key format."
    fi

    log "Fetched CA public key: $TRUSTSSH_CA_PUBLIC_KEY"
}

download_bootstrap_tooling() {
    log "Downloading TrustSSH bootstrap tooling"

    rm -rf "$WORKDIR"
    mkdir -p "$WORKDIR"

    git clone --depth 1 --branch "$TRUSTSSH_BRANCH" "$TRUSTSSH_REPO" "$WORKDIR/TrustSSH"
}

run_ansible_local() {
    log "Running TrustSSH Ansible bootstrap locally"

    cd "$BOOTSTRAP_DIR"

    cat > /tmp/trustssh-local-inventory.ini <<'EOF'
[trustssh_servers]
localhost ansible_connection=local
EOF

    ansible-playbook playbooks/install-trustssh-server.yml \
        -i /tmp/trustssh-local-inventory.ini \
        -e "{\"trustssh_ca_public_key\":\"$TRUSTSSH_CA_PUBLIC_KEY\"}"
}

main() {
    require_root
    validate_endpoint
    detect_package_manager
    install_packages
    fetch_public_key
    download_bootstrap_tooling
    run_ansible_local

    log "TrustSSH server bootstrap complete"
}

main "$@"