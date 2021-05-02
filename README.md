# TFselect

Inspired by  [warrensbox/terraform-switcher](https://github.com/warrensbox/terraform-switcher)

The `tfselect` command line tool lets you switch between different versions of [terraform](https://www.terraform.io/).
If you do not have a particular version of terraform installed, `tfselect` will download the version you desire.
The installation is minimal and easy.
Once installed, simply select the version you require from the dropdown and start using terraform.

## Installation

`tfselect` is available for MacOS and Linux based operating systems (Windows experemetal).

### Homebrew

Installation for MacOS/Linux is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/).

```ruby
brew install basti0nz/tap/tfselect
```

### NPM

```sh
npm i @versusdev/tfselect
```

### Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/basti0nz/tfselect/release/install.sh | bash
```

### Build and install SNAP package

```bash
snap install snapcraft --classic

snap install multipass

snapcraft

snap install tfselect_*.snap --devmode --dangerous

tfselect -v

multipass stop snapcraft-tfselect && multipass delete snapcraft-tfselect && multipass purge

```

### Get binary releases or install from source

Alternatively, you can get releases or install the binary from source [here](https://github.com/basti0nz/tfselect/releases)

## How to use

### Use dropdown menu to select version

1. You can switch between different versions of terraform by typing the command `tfselect` on your terminal.
2. Select the version of terraform you require by using the up and down arrow.
3. Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfselect 0.13.5` for version 0.10.5 of terraform.
3. Hit **Enter** to switch.

### Install latest version only

1. Install the latest stable version only.
2. Run `tfselect -u` or `tfselect --latest`.
3. Hit **Enter** to install.

### Install latest implicit version for stable releases

1. Install the latest implicit stable version.
2. Ex:  `tfselect 0.13` downloads 0.13.6 (latest) version.
3. Hit **Enter** to install.

### Use version.tf file  

If a version.tf file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the lastest version:

```ruby
terraform {
  required_version = ">= 0.12.9"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}
```

### Use .tfselect.toml file  (For non-admin - users with limited privilege on their computers)

This is similiar to using a .tfswitchrc file, but you can specify a custom binary path for your terraform installation

1. Create a custom binary path. Ex: `mkdir /Users/basti0nz/bin` (replace basti0nz with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/basti0nz/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terraform. Ex: `tfselect -b /Users/basti0nz/bin/terraform 0.10.8 `
4. Optionally, you can create a `.tfselect.toml` file in your home directory for global settings.
5. Your `.tfselect.toml` file should look like this:

```ruby
bin = "/Users/versus/bin/terraform"
version = "0.11.3"
```

6. Run `tfselect` and it should automatically install the required terraform version in the specified binary path

Alternatively, you can generate .tfselect.toml in current directory just use `tfselect --init ` or `tfselect --init 0.13.2`

### Use .tfswitchrc file

1. Create a `.tfswitchrc` file containing the desired version
2. For example, `echo "0.10.5" >> .tfswitchrc` for version 0.10.5 of terraform
3. Run the command `tfselect` in the same directory as your `.tfswitchrc`

### Use environment variable TFSELECT_PATH

1. Create a `TFSELECT_PATH` environment variable with your custom path to install terraform Ex: `export  TFSELECT_PATH=/Users/versus/bin/terraform`
2. Run the command `tfselect` 

### Use environment variable TFSELECT_VERSION

1. Create a `TFSELECT_PATH` environment variable with your custom path to install terraform Ex: `export  TFSELECT_VERSION=0.13.5`
2. Run the command `tfselect`

#### *File `.terraform-version` may be used for compatibility with [`tfenv`](https://github.com/tfutils/tfenv#terraform-version-file) and other tools which use it*

### Use terragrunt.hcl file

If a terragrunt.hcl file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the lastest version:

```ruby
terragrunt_version_constraint = ">= 0.26, < 0.27"
terraform_version_constraint  = ">= 0.13, < 0.14"
...
```

#### Example use case on a CI/CD agent

An agent host is scheduled to run multiple Terraform Plan/Apply tasks simultaneously from different repositories, which each contain different Terraform configurations, using different versions of Terraform.

With a regular symlink approach, they would have to be queued and run synchronously, as a version switch would have to be made between each run.

However, if the terraform binary (or symlink) is replaced with a script (example below) that catches the arguments sent to terraform, performs a tfselect using both the no-symlink and quiet flags, then catches the output from tfselect and uses the desired binary to run the specific task, this can be done dynamically with each run, in parallel.

This utilizes tfselect's brilliant ability to detect and install the required version, while allowing concurrency.

Example PoC script

```sh
#!/bin/bash

TF_DEFAULT_VERSION=0.13.2
TF_BIN_LOCATION=~/.terraform.versions

TF_SELECT=`tfselect -nq`
TF_DETECTED_VERSION=`echo "$TF_SELECT" | awk -F '\"|\"' '{print $2}' | grep "\S"`

echo "$TF_SELECT"

if [ "$TF_DETECTED_VERSION" == "" ]; then
    TF_DETECTED_VERSION=$TF_DEFAULT_VERSION
fi

TF_BIN="$TF_BIN_LOCATION/terraform_$TF_DETECTED_VERSION"

echo "Using version: $TF_DETECTED_VERSION"
echo ""

$TF_BIN "$@"
```

Suggestions and improvements are welcome.

## Issues

Please open  *issues* here: [New Issue](https://github.com/basti0nz/tfselect/issues)
