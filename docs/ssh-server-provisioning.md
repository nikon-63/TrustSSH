# SSH Server Provisioning

## Overview

There are three main ways to provision a Linux SSH server to trust the TrustSSH CA:

1. **Automatic bootstrap installer** - Recommended for quickly provisioning a single server by running the installer directly on the SSH server.
2. **Manual Ansible playbook** - Recommended when provisioning multiple servers.
3. **Fully manual configuration** - If you just want to do it yourself.

## Automatic Bootstrap Installer

The easiest way to provision a server is to run the bootstrap installer directly on the SSH server.

Set your TrustSSH endpoint first:

```bash
export TRUSTSSH_ENDPOINT="https://trustssh.example.com"
```

Then run the installer:

```bash
curl --proto '=https' --tlsv1.2 -sSfL https://raw.githubusercontent.com/nikon-63/TrustSSH/main/helpers/bootstrap-install.sh | sudo -E bash
```

## Manual Ansible Provisioning

The `bootstrap-tooling` directory contains Ansible playbooks for installing and removing the TrustSSH server configuration.

### Installation

From inside the `bootstrap-tooling` directory, run:

```bash
ansible-playbook -i inventory.ini playbooks/install-trustssh-server.yml \
  -e '{"trustssh_ca_public_key":"ssh-ed25519 AAAA... trustssh-ca"}'
```

### Remove TrustSSH CA Trust

To remove the TrustSSH CA public key and OpenSSH configuration:

```bash
ansible-playbook -i inventory.ini playbooks/remove-trustssh-server.yml
```

## Fully Manual Provisioning

Manual provisioning can be used when you want to configure the server yourself.

1. Upload the TrustSSH CA public key to the SSH server. 

```bash
sudo tee /etc/ssh/trustssh_ca.pem >/dev/null <<'EOF'
ssh-ed25519 AAAA... trustssh-ca
EOF
```


2. (a) On systems that support OpenSSH drop-in configuration, create:

```bash
sudo tee /etc/ssh/sshd_config.d/99-trustssh.conf >/dev/null <<'EOF'
TrustedUserCAKeys /etc/ssh/trustssh_ca.pem
EOF
```

2. (b) If the system does not use `/etc/ssh/sshd_config.d/*.conf`, add the following line directly to `/etc/ssh/sshd_config`:

```text
TrustedUserCAKeys /etc/ssh/trustssh_ca.pem
```

3. Validate the SSH server configuration:

```bash
sudo sshd -t
```

4. Restart the SSH server:

```bash
sudo systemctl restart sshd
```

## Notes

- The automatic bootstrap installer expects the TrustSSH endpoint to expose the CA public key at `/public_key.txt`.
- The CA public key should be a public key only. Do not place the TrustSSH CA private key on SSH servers.
- On minimal Ubuntu containers, `/run/sshd` may need to exist before `sshd -t` succeeds.
- Alpine Linux usually requires Python to be installed before Ansible can run.
