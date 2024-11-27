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

# Fix inside env. cobra-cli and go path

To ensure that the cobra-cli command remains available even after restarting your PC, you need to add the path to the executable to your system's PATH environment variable permanently. The command export PATH=$PATH:$GOBIN you're using only sets the PATH for the current terminal session. Once the terminal is closed or the PC is restarted, this setting is lost.

To make this change permanent, you need to add the export command to a shell startup file like .bashrc, .bash_profile, or .zshrc, depending on your shell and operating system.

Here's how you can do it:

Open your shell's startup file: This file is typically located in your home directory. If you are using Bash, it’s usually .bashrc or .bash_profile. If you're using Zsh (common in newer macOS versions), it’s .zshrc.

For Bash on Linux or macOS:
```bash
nano ~/.bashrc
```

For Bash on macOS (if you use .bash_profile):
```bash
nano ~/.bash_profile
```

For Zsh (common on newer macOS versions):
```bash
nano ~/.zshrc
```

Add the export command to the file: At the end of the file, add the following line:
```bash
export GOBIN=~/go/bin
export PATH=$PATH:$GOBIN
```

Replace $GOBIN with the actual path to your Go binaries if it's not already set as an environment variable (you can do it using: "echo $(go env GOPATH)/bin").

Save and close the file: If you’re using nano, you can do this by pressing Ctrl + O, Enter, and then Ctrl + X.

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


# Cobra cli usage

- Only to initialize the repo go mod init github.com/eltiocaballoloco/sinaloa-cli
- go install github.com/spf13/cobra-cli@latest
- go get -u github.com/ricochet2200/go-disk-usage/du
- go get -u github.com/ricochet2200/go-disk-usage
- go get -u storj.io/uplink
- go get github.com/stretchr/testify/assert
- go get github.com/stretchr/testify/mock
- go get gopkg.in/yaml.v2
- cobra-cli init <cli_name> (to create cli)
- cobra-cli add <name_cmd> (to create a new command)


# Development

- Install required
- Install deps


# Debug usage



# Makefile



# Env

It is necessary create an .env file:

```bash
SINALOA_CLI_DEBUG="false" # true or false
STORJ_SECRET="YOUR_SECRET" # Storj secret
```

Now you can do:

```bash
cd scripts
```
```bash
source set_env_var.sh
```

So you are able to set OS ENV variables before the execution of the sinaloa cli.