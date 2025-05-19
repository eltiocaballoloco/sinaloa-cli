# SINALOA-CLI
The devops CLI used for automations.


# Go installation

- MacOS
```bash
brew install golang
export GOBIN=~/go/bin
export PATH=$PATH:$GOBIN
```

- Ubuntu
```bash
sudo apt install golang
export GOBIN=~/go/bin
export PATH=$PATH:$GOBIN
```

# Fix inside env cobra-cli and go path

To ensure that the cobra-cli command remains available even after restarting your PC, you need to add the path to the executable to your system's PATH environment variable permanently. The command export PATH=$PATH:$GOBIN you're using only sets the PATH for the current terminal session. Once the terminal is closed or the PC is restarted, this setting is lost.

To make this change permanent, you need to add the export command to a shell startup file like .bashrc, .bash_profile, or .zshrc, depending on your shell and operating system.

Here's how you can do it:

Open your shell's startup file: This file is typically located in your home directory. If you are using Bash, it’s usually .bashrc or .bash_profile. If you're using Zsh (common in newer macOS versions), it’s .zshrc.

Add the path of the go-bin to make available cobra-cli on the os:
```bash
echo $(go env GOPATH)/bin
/eltio/go/bin --> path of go bin
```
Then go to the .zshrc or .bashrc file, save and close the file.

For Linux:
```bash
nano ~/.bashrc
```
or
```bash
nano ~/.bash_profile
```
For Zsh (common on newer macOS versions):
```bash
nano ~/.zshrc
```

Apply the changes: For the changes to take effect, you need to reload the startup file. You can do this by either restarting your terminal or sourcing the file with one of the following commands:

For .bashrc or .bash_profile:
```bash
source ~/.bashrc
# or
source ~/.bash_profile
```

For .zshrc:
```bash
source ~/.zshrc
```

After doing this, the cobra-cli command should be available in all new terminal sessions. This way, the PATH update becomes a permanent part of your shell configuration.


# Packages and cobra cli usage

- Only to initialize the repo go mod init github.com/eltiocaballoloco/sinaloa-cli
- go install github.com/spf13/cobra-cli@latest
- go get github.com/stretchr/testify/assert
- go get github.com/stretchr/testify/mock
- go get gopkg.in/yaml.v2
- go get github.com/microsoftgraph/msgraph-sdk-go@v1.53.0
- go get github.com/microsoft/kiota-http-go@v1.4.5
- go get github.com/Azure/azure-sdk-for-go/sdk/azidentity@v1.8.0
- cobra-cli init <cli_name> (to create cli)
- cobra-cli add <name_cmd> (to create a new command)


# Development

- Install required
- Install deps


# Makefile

This project includes a `Makefile` to simplify common tasks for building, testing, and managing the `sinaloa-cli` application.

## Targets

### `make build`
- Compiles the Go source code from `src/main.go`.
- Places the output binary in the `build/` directory.
- Copies the final binary to `/usr/local/bin/sinaloa` for global CLI use.

---

### `make build-clean`
- Cleans Go build artifacts.
- Deletes the `build/` directory.

---

### `make deps`
- Installs all Go dependencies.
- Equivalent to:  
  ```bash
  go get -v -t -d ./...
  ```

---

### `make mod-tidy`
- Cleans up `go.mod` and `go.sum` using:
  ```bash
  go mod tidy
  ```

---

### `make test`
- Runs all Go tests located in the `tests/` directory.

---

## Command Generation

### `make new-cmd cmd=<command> subcmd=<subcommand> flags="<flags>"`
Generates a new command and subcommand with the specified flags using `scripts/create_cmd.sh`.

- Example:
  ```bash
  make new-cmd cmd=storj subcmd=add flags="secret:secret:s:Storj secret to connect within:true:|path:path:p:Path to store file:true:"
  ```

---

### `make new-sub cmd=<command> subcmd=<subcommand> flags="<flags>"`
Adds a new subcommand to an existing command using `scripts/create_sub.sh`.

- Example:
  ```bash
  make new-sub cmd=storj subcmd=put flags="msg:msg:m:Message to receive:true:|path:path:p:Storage path:true:"
  ```

---

## Notes
- Ensure that `scripts/create_cmd.sh` and `scripts/create_sub.sh` are executable.
- Requires Go installed and accessible via the `go` command.
- Run `make build` with `sudo` if you want to copy the binary to `/usr/local/bin`.


# Env

It is necessary create an .env file:

```bash
AZURE_TENANT_ID="xxx-yyyy-tttt-1234"
AZURE_CLIENT_ID="xxx-yyyy-tttt-1234"
AZURE_CLIENT_SECRET="xxx-yyyy-tttt-1234"
AZURE_DRIVE_ID="xxx-yyyy-tttt-1234"
```

Now you can do:

```bash
cd scripts
```
```bash
source set_env_var.sh
```

So you are able to set OS ENV variables before the execution of the sinaloa cli.