package letsgo

import (
	"fmt"
	"log"
	"os"
	"github.com/spf13/cobra"
	"regexp"
	"os/user"
	"io"
	"io/ioutil"
	"strings"
)

// Variables
var cwd string
var gopath string
var path string
var mkdir bool = false
var zshrc bool = false
var bin bool = false
var src bool = false
var pkg bool = false
var match_gopath, match_path bool
var export_cmds string

// const colors
const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

// Core Logic
func updateEnv() (path string, path_env string, export_cmds string) {
	/* This function will update the GOPATH and PATH with the current directory path
	*/

	// Collect crrent path
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Check for bin src and pkg dirs
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	dir_list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir_list {
		if file.Name() == "bin" {
			bin = true
		}
		if file.Name() == "src" {
			src = true
		}
		if file.Name() == "pkg" {
			pkg = true
		}
	}

	// Display commands for GOPATH and PATH
	if !(bin && pkg && src) {
		fmt.Println(string(colorRed), "One or all of bin/, src/ and pkg/ directories not present!!", string(colorReset))
		path_env = path
	} else {
		fmt.Println(string(colorGreen), "All of bin/, src/ and pkg/ directories are present..", string(colorReset))
		path_env = path + "/bin"
	} 

	match, err := regexp.MatchString(path, os.Getenv("PATH"))
	if err != nil {
		log.Fatal(err)
	}
	if match {
		fmt.Println("current dir already in $PATH")
		export_cmds = "export GOPATH="+path+"\n"
	} else {
		export_cmds = "export GOPATH="+path+"\n"
		export_cmds += "export PATH=$PATH:"+path_env+"\n"
	}
	return path, path_env, export_cmds
}

func dirSetup() {
	// create directories in local path]
	err := os.Mkdir("bin", 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir("pkg", 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir("src", 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(colorGreen), "Go directories created...", string(colorReset))
}

func updateZshrc(gopath string, path string)  {
	// Update ~/.zshrc file
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	//Create backup zsh
	zsh, err := os.Open(user.HomeDir+"/.zshrc")
	if err != nil {
		log.Fatal(err)
	}
	defer zsh.Close()
	
	zsh_bak, err := os.Create(user.HomeDir+"/.zshrc.bak")
	if err != nil {
		log.Fatal(err)
	}
	defer zsh_bak.Close()
	_, err = io.Copy(zsh_bak, zsh)
	if err != nil {
		log.Fatal(err)
	}
	defer zsh_bak.Close()

	// Update file
	zsh_data, err := ioutil.ReadFile(user.HomeDir+"/.zshrc")
	if err != nil {
		log.Fatal(err)
	}
	zsh_temp := strings.Split(string(zsh_data), "\n")

	for index, item := range zsh_temp {
		check_gopath, err := regexp.MatchString("export GOPATH=", item)
		if err != nil {
			log.Fatal(err)
		}
		check_path, err := regexp.MatchString("export PATH=\\$PATH:/Users/rjayaraman/terminal_files/repos/letsgo/bin", item)
		if err != nil {
			log.Fatal(err)
		}
		if check_gopath {
			zsh_temp[index] = "export GOPATH="+gopath
			match_gopath = true
		}
		if check_path {
			zsh_temp[index] = "export PATH=$PATH:"+path
			match_path = true
		}
	}
	// Update zshrc in place
	if match_gopath || match_path {
		zsh_update := strings.Join(zsh_temp, "\n")
		err := ioutil.WriteFile(user.HomeDir+"/.zshrc", []byte(zsh_update), 0644)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(colorGreen), "Updated zshrc file in place", string(colorReset))
		}
	}

	// Append to zshrc
	zsh_file, err := os.OpenFile(user.HomeDir+"/.zshrc", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer zsh_file.Close()
	if !match_gopath {
		if _, err = zsh_file.WriteString("\nexport GOPATH="+gopath); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(colorGreen), "Appended GOPATH to zshrc", string(colorReset))
		}
	}
	if !match_path {
		if _, err = zsh_file.WriteString("\nexport PATH=$PATH:"+path); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(colorGreen), "Appended PATH to zshrc", string(colorReset))
		}
	}
}

// COBRA CLI
var rootCmd = &cobra.Command{
	Use:   "letsgo",
	Short: "CLI app that prepares the current working directory for go development",
	Long: `This cli app does the following
	1. Create the bin/, src/ and pkg/ directories if not present in current directory if flag is present
	2. Update GOPATH and PATH
	3. Update the zshrc file if flag is present`,

	// Core logic
	Run: func(cmd *cobra.Command, args []string) {
		if mkdir {
			dirSetup()
		}
		gopath, path, export_cmds = updateEnv()
		if zshrc {
			updateZshrc(gopath, path)
			export_cmds = "sorurce ~/.zshrc"
		}
		fmt.Println("\n\nRun the following commands:-")
		fmt.Println(export_cmds)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&mkdir, "mkdir", "m", false, "crate bin/, src/, pkg/ directories")
	rootCmd.Flags().BoolVarP(&zshrc, "zshrc", "z", false, "update the ~/.zshrc file with new GOPATH and PATH")
}
