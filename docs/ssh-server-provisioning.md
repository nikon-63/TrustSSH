# SSH Server Provisioning

## Overview
There are three main ways to provision a Linux SSH server to trust the TrustSSH CA:

1. **Fully automated Ansible playbook** - Recommended approach but installs a lot of tooling on the servers.
2. **Manual Ansible playbook** - Server provisioning controlled by ansible playbook but run manually on your mac.
3. **Fully manual** - Manually perform edit the SSH server configuration files and restart the SSH service.

### Automated Ansible Playbook
TODO

### Manual Ansible Playbook
With in the bootstrap-tooling dir there are two ansible playbooks. The `install-trustssh-server.yml` playbook will install the TrustSSH CA public key on the server and configure SSH to trust it. The `remove-trustssh-server.yml` playbook will remove the TrustSSH CA public key and configuration.

To run the playbook you will need to create an inventory file from the example below and then run the playbook using `ansible-playbook` command.

```bash
# From with in the bootstrap-tooling directory
# Install TrustSSH CA public key on the server
ansible-playbook playbooks/install-trustssh-server.yml \
 -e '{"trustssh_ca_public_key":"ssh-ed25519 AAAA BBB CCC trustssh-ca"}'
```

```bash
# From with in the bootstrap-tooling directory
# Remove TrustSSH CA public key from the server
ansible-playbook playbooks/remove-trustssh-server.yml
```

### Fully Manual Provisioning
You can manually configure your SSH server to trust the TrustSSH CA by uploading and adding the CA public key to the `TrustedUserCAKeys` in your `sshd_config` file.

```bash
# /etc/ssh/sshd_config
TrustedUserCAKeys /etc/ssh/trustssh_ca.pem
```