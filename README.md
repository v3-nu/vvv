# VV - A CLI Utility for various tasks
VV is a CLI utility that was written to assist with various CLI tasks. It contains a number of commands designed to make a bunch of common tasks easier.

## Installation
### Package Managers
You can install VV using a package manager.

#### Debian/Ubuntu
```
# Download and save the repository key
sudo curl https://pkg.cloudyne.io/debian/repository.key -o /etc/apt/trusted.gpg.d/cloudyne.asc

# Add the repository to your sources list
echo "deb https://pkg.cloudyne.io/debian/ stable main" | sudo tee /etc/apt/sources.list.d/cloudyne.list

# Update your package list
sudo apt update

# Install VV
sudo apt install vv
```
### Go Install
You can also install VV using Go.

```
go install github.com/cloudyne-io/vv
```

### Binaries
You can also download the binaries from the [Releases](https://github.com/v3-nu/vvv/releases) page.

## Commands
- [x] vv
  - [x] alias - `Add persistent aliases to your shell. The aliases are stored in a json file, and easily editable or movable.`
    - [x] list - `List all aliases`
    - [x] set - `Set a new value for an alias`
    - [ ] install - `Install the alias configuration into your shell`
  - [x] cryptography
    - [x] pfx-to-pem (pfx2pem, p2p) - `Convert a pfx file to a set of PEM files (.key, .nopass.key, .crt, .ca.crt, .fullchain.pem)`
    - [ ] pem-to-pfx (p2pfx, p2p) - `Convert one or multiple PEM files to a pfx file`
  - [x] install (i)
    - [x] azure-cli (azurecli) - `Install Azure CLI (Debian)`
    - [x] go - `Install Go`
    - [x] helm - `Install Helm`
    - [x] kubectl - `Install Kubectl`
    - [x] jupyter-bash-kernel - `Install Jupyter Bash Kernel`
    - [x] nodejs - `Install NodeJS`
    - [x] acme-shell - `Install acme.sh`
  - [ ] kubectl (k, k8s)
    - 


## Default Aliases
- [ ] vvk - vvv vv kubectl
- [ ] vvi - vvv install
- [ ] vvp - vvv package
- [ ] kubens - vvv kubectl set-ns
- [ ] kubectx - vvv kubectl getctx
- [ ] kubeuse - vvv kubectl use-context
- [ ] kubesec - vvv kubectl getsecret