# CLI Installation Using Homebrew

> [!NOTE]  
> After installing the CLI using Homebrew, you can skip the [Build the Binary](./cli-deployment.md#build-the-binary) section of the deployment guide and go to [Configure the CLI](./cli-deployment.md#configure-the-cli).

TrustSSH can be installed on macOS and Linux using Homebrew. This is the recommended installation method because it downloads the latest published release and makes the `trustssh` command available system-wide.

## Quick install

Use this one-line command to install TrustSSH:

```bash
brew install nikon-63/tap/trustssh
```

This automatically adds the `nikon-63/tap` Homebrew tap and installs the latest TrustSSH release.

## Manual install

You can also add the tap first and then install TrustSSH separately:

```bash
brew tap nikon-63/tap
brew install trustssh
```

## Verify the installation

After installation, check that the CLI is available:

```bash
which trustssh
trustssh
```

You should see the installed TrustSSH binary path and the CLI usage output.

## Update TrustSSH

To update Homebrew and upgrade TrustSSH to the latest release:

```bash
brew update
brew upgrade trustssh
```

## Uninstall TrustSSH

To remove the TrustSSH CLI:

```bash
brew uninstall trustssh
```

To also remove the Homebrew tap:

```bash
brew untap nikon-63/tap
```
