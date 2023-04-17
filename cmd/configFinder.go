package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/basti0nz/tfselect/lib"
	"github.com/kiranjthomas/terraform-config-inspect/tfconfig"
	"github.com/spf13/viper"
)

const (
	tfvFilename  = ".terraform-version"
	rcFilename   = ".tfswitchrc"
	tomlFilename = ".tfselect.toml"
	//TODO: add support .tfswitch.toml
	tgFilename = "terragrunt.hcl"
	defaultBin = "/usr/local/bin/terraform" //default bin installation dir
)

func GetConfigVariable() (string, string) {
	var path string
	exist, tfversion, path := findConfig()
	if exist == true {
		if path == "" {
			if exist, _, p := checkHomeDirToml(); exist == true {
				return tfversion, p
			}
			return tfversion, defaultBin
		}
		return tfversion, path
	}

	if IsUserRoot() {
		path = defaultBin
	} else {
		dirname, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if _, err := os.Stat(dirname + "/.bin/"); !os.IsNotExist(err) {
			path = dirname + "/.bin/terraform"
		} else if _, err := os.Stat(dirname + "/bin/"); !os.IsNotExist(err) {
			path = dirname + "/bin/terraform"
		} else {
			fmt.Println("Error: Cannot find terraform binary path include `bin` and `.bin` in your home directory. Use flag -b to custom binary path.")
			os.Exit(1)
		}
	}
	return "", path
}

func GetInstalledVersion(bin string) (string, error) {
	cmd := exec.Command(bin, "-v")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	version := trimFirstRune(strings.Split(outb.String(), " ")[1])
	return strings.Split(strings.TrimSpace(version), "\n")[0], nil
}

func findConfig() (bool, string, string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	var tfversion, path string
	var exist bool

	if exist, tfversion, path = checkToml(dir); exist == true {
		return exist, tfversion, path
	}

	if exist, tfversion = checkTGfile(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion = checkTFVAR(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion = checkTFswitchrc(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion, path = checkHomeDirToml(); exist == true {
		return exist, tfversion, path
	}

	return false, "", ""
}

// TODO: refactoring for support  .tfswitch.toml and tfselect.toml
func checkToml(dir string) (bool, string, string) {
	configfile := dir + fmt.Sprintf("/%s", tomlFilename) //settings for .tfselect.toml file in current directory (option to specify bin directory)
	if _, err := os.Stat(configfile); err == nil {
		fmt.Printf("Reading configuration from %s\n", configfile)

		var path string                                 //takes the default bin (defaultBin) if user does not specify bin path
		configfileName := lib.GetFileName(tomlFilename) //get the config file
		viper.SetConfigType("toml")
		viper.SetConfigName(configfileName)
		viper.AddConfigPath(dir)

		errs := viper.ReadInConfig() // Find and read the config file
		if errs != nil {
			fmt.Printf("Unable to read %s provided\n", tomlFilename) // Handle errors reading the config file
			fmt.Println(err)
			os.Exit(1) // exit immediately if config file provided but it is unable to read it
		}

		bin := viper.Get("bin") // read custom binary location
		path = os.ExpandEnv(bin.(string))
		tfversion := viper.Get("version") //attempt to get the version if it's provided in the toml
		return true, tfversion.(string), path
	}

	return false, "", ""
}

func checkHomeDirToml() (bool, string, string) {
	usr, errCurr := user.Current()
	if errCurr != nil {
		return false, "", ""
	}
	return checkToml(usr.HomeDir)
}

func checkTFswitchrc(dir string) (bool, string) {
	rcfile := dir + fmt.Sprintf("/%s", rcFilename) //settings for .tfswitchrc file in current directory (backward compatible purpose)
	if _, err := os.Stat(rcfile); err == nil {
		fmt.Printf("Reading required terraform version %s \n", rcFilename)
		fileContents, err := ioutil.ReadFile(rcfile)
		if err != nil {
			fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/versus/terraform-switcher/blob/master/README.md\n", rcFilename)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		tfversion := strings.TrimSuffix(string(fileContents), "\n")
		return true, tfversion
	}
	return false, ""
}

func checkTFvfile(dir string) (bool, string) {
	tfvfile := dir + fmt.Sprintf("/%s", tfvFilename) //settings for .terraform-version file in current directory (tfenv compatible)

	if _, err := os.Stat(tfvfile); err == nil { //if there is a .terraform-version file, and no command line arguments
		fmt.Printf("Reading required terraform version from %s \n", tfvFilename)

		fileContents, err := ioutil.ReadFile(tfvfile)
		if err != nil {
			fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/versus/terraform-switcher/blob/master/README.md\n", tfvFilename)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		tfversion := strings.TrimSuffix(string(fileContents), "\n")
		return true, tfversion
	}
	return false, ""
}

func checkTGfile(dir string) (bool, string) {
	var tfconstraint string
	tgfile := dir + fmt.Sprintf("/%s", tgFilename)
	if _, err := os.Stat(tgfile); err == nil { //if there is a terragrunt file, and no command line arguments
		fmt.Printf("Reading required terraform version from %s \n", tgFilename)

		fileContents, err := ioutil.ReadFile(tgfile)
		if err != nil {
			fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/versus/terraform-switcher/blob/master/README.md\n", tgFilename)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		regex, _ := regexp.Compile(`^terraform_version_constraint\s+=\s+"(?P<version>.*)".*`)
		for _, line := range strings.Split(strings.TrimSuffix(string(fileContents), "\n"), "\n") {
			if regex.MatchString(line) {
				res := regex.FindStringSubmatch(line)
				tfconstraint = res[1]
				break
			}
		}

		return getFromConstraint(tfconstraint)
	}
	return false, ""
}

func checkTFVAR(dir string) (bool, string) {
	if module, _ := tfconfig.LoadModule(dir); len(module.RequiredCore) >= 1 { //if there is a version.tf file, and no commmand line arguments
		fmt.Printf("Reading required terraform version from version.tf \n")
		tfconstraint := module.RequiredCore[0] //we skip duplicated definitions and use only first one
		return getFromConstraint(tfconstraint)
	}
	return false, ""
}

func getFromConstraint(tfconstraint string) (bool, string) {
	var tfversion string
	constrains, err := semver.NewConstraint(tfconstraint) //NewConstraint returns a Constraints instance that a Version instance can be checked against
	if err != nil {
		return false, ""
	}
	listAll := true

	tflist, _ := lib.GetTFList(hashiURL, listAll)
	versions := make([]*semver.Version, len(tflist))
	for i, tfvals := range tflist {
		version, err := semver.NewVersion(tfvals) //NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
		if err != nil {
			return false, ""
		}
		versions[i] = version
	}
	sort.Sort(sort.Reverse(semver.Collection(versions)))
	for _, element := range versions {
		if constrains.Check(element) { // Validate a version against a constraint
			tfversion = element.String()
			fmt.Printf("Matched version: %s\n", tfversion)
			return true, tfversion
		}
	}
	return false, ""
}

func trimFirstRune(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return ""
}
