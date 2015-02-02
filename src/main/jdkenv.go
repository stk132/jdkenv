package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	jdkdir        = homeDir() + "/.jdkenv/java"
	macSystemJdk  = "/System/Library/Java/JavaVirtualMachines/"
	macLibraryJdk = "/Library/Java/JavaVirtualMachines/"
)

func main() {
	//var jdkdir = homeDir() + "/.jdkenv/java"
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("print help")
	} else {
		switch args[0] {
		case "init":
			initialize()
		case "list":
			list()
		case "use":
			if len(args) >= 2 {
				use(args[1])
				return
			} else {
				fmt.Println("please select jdk directory")
			}
		case "current":
			current()
		default:
			fmt.Println("print help")
		}
	}

	os.Exit(0)
}

func use(ver string) {
	if runtime.GOOS == "darwin" {
		macUse(ver)
		return
	}
	if !exist(jdkdir + "/" + ver) {
		fmt.Println(ver + "is not exist")
		return
	}

	jdkpath := jdkdir + "/" + ver
	javahomesymlink := jdkdir + "/current"

	removeCurrnetSymlink(javahomesymlink)
	makeJavahomeSymlink(jdkpath, javahomesymlink)
}

func macUse(ver string) {
	var jdkpath string
	if exist(macSystemJdk + ver) {
		jdkpath = macSystemJdk + ver + "/Contents/Home"
	} else if exist(macLibraryJdk + ver) {
		jdkpath = macLibraryJdk + ver + "/Contents/Home"
	} else {
		fmt.Println(ver + " isn't exists at this System")
		return
	}

	javahomesymlink := jdkdir + "/current"

	removeCurrnetSymlink(javahomesymlink)
	makeJavahomeSymlink(jdkpath, javahomesymlink)
}

func removeCurrnetSymlink(javahomesymlink string) {
	if exist(javahomesymlink) {
		if err := os.Remove(javahomesymlink); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func makeJavahomeSymlink(jdkpath, javahomesymlink string) {
	if err := os.Symlink(jdkpath, javahomesymlink); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func list() {
	if runtime.GOOS == "darwin" {
		macJdkList()
		return
	}
	dirs, err := ioutil.ReadDir(jdkdir)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(dirs) == 0 {
		fmt.Println("jdk isn't exists at " + jdkdir)
		return
	}

	for _, value := range dirs {
		if strings.HasPrefix(value.Name(), "jdk") {
			fmt.Println(value.Name())
		}
	}
}

func macJdkList() {
	printMacJdk(macSystemJdk)
	printMacJdk(macLibraryJdk)
}

func printMacJdk(dirPath string) {
	dirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, value := range dirs {
		fmt.Println(value.Name())
	}
}

func current() {
	javahomesymlink := jdkdir + "/current"
	if !exist(javahomesymlink) {
		fmt.Println("jdkenv not used")
		return
	}

	dest, err := os.Readlink(javahomesymlink)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch runtime.GOOS {
	case "darwin":
		{
			splitedpath := strings.Split(dest, string(os.PathSeparator))
			fmt.Println(splitedpath[len(splitedpath)-3])
		}
	default:
		fmt.Println(filepath.Base(dest))
	}

}

func initialize() {
	if !exist(jdkdir) {
		err := os.MkdirAll(jdkdir, 0777)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	switch runtime.GOOS {
	case "windows":
		windowsInit()
	default:
		unixTypeInit()
	}
}

func windowsInit() {
	_, err := exec.Command("setx", "JAVA_HOME", jdkdir+"/current").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("please reboot command prompt to recognize JAVA_HOME")

	if hasGitBash() {
		fmt.Println("if you use git bash, write in your .bashrc below")
		printSetJavaHomeMsg()
	}
}

func unixTypeInit() {
	fmt.Println("write in your .bashrc below")
	printSetJavaHomeMsg()
}

func printSetJavaHomeMsg() {
	fmt.Println("export JAVA_HOME=" + jdkdir + "/current")
	fmt.Println("and execute below")
	fmt.Println(". ~/.bashrc")
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return usr.HomeDir
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasGitBash() bool {
	return runtime.GOOS == "windows" && os.Getenv("HOME") == homeDir()
}
