package main

import (
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"bufio"
	"path/filepath"
	"strings"
	"errors"
)


const (
	Help = `First you need to configure this command, to configure use:
	rcp --configure <secret> <host>
	where:
	<secret> must be replaced by a string that only you knowns and it's hard to guess, recommended some hash like https://passwordsgenerator.net/sha256-hash-generator/
	<host> is optional, if you don't known if you should put it so don't, default is: https://rcp-cv.herokuapp.com

	
Args
rcp --configure | -c | --copy | -p | --paste | -e | --erase [text]

-c --copy -> copy [text]
-p --pase -> prints [text] previous copied
-e --erase -> clears any previous copied [text]

examples:
	rcp -c hello world!
	echo "hello world!" | rcp -c
	rcp -p | echo`

	ConfigFileContent = `secret <secret>
host <host>
`

	DefaultHost = "https://rcp-cv.herokuapp.com"
	ConfigFileName = ".rcp.properties"
)

type Conf struct {
	Secret, Host string
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println(Help)
		os.Exit(1)
	}

	if os.Args[1] == "--configure" {
		if len(os.Args) < 3 {
			fmt.Println(Help)
			os.Exit(1)
		}
		host := DefaultHost
		if len(os.Args) == 4 {
			host = os.Args[3]
		}

		err := configure(os.Args[2], host)
		checkFailAndExit(err)
		os.Exit(0)
	}
	a := os.Args[1]

	if a == "-c" || a == "--copy" {
		c, err := recoverConfig()
		checkFailAndExit(err)
		copyy(getText(), c)
	} else if a == "-p" || a == "--paste" {
		c, err := recoverConfig()
		checkFailAndExit(err)
		err = paste(c)
		checkFailAndExit(err)
	} else if a == "-e" || a == "--erase" {
		c, err := recoverConfig()
		checkFailAndExit(err)
		err = erase(c)
		checkFailAndExit(err)
	} else if a == "-h" || a == "--help" {
		fmt.Println(Help)
		os.Exit(0)
	} else {
		fmt.Println(Help)
		os.Exit(1)
	}
}

func getText() string {
	fileInfo, _ := os.Stdin.Stat()
	isPipe := fileInfo.Mode() & os.ModeCharDevice == 0
	if isPipe {
		txt := ""
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			txt = txt + scanner.Text()
		}
		return txt
	} else {
		return strings.Join(os.Args[2:], " ")
	}
}

func checkFailAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func recoverConfig() (Conf, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Conf{}, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Conf{}, err
	}

	secret := ""
	host := ""

	ll := strings.Split(string(data), "\n")
	for _, l := range ll {
		props := strings.Split(l, " ")
		if props[0] == "#" {
			continue
		} else if props[0] == "secret" {
			secret = props[1]
		} else if props[0] == "host" {
			host = props[1]
		}
	}

	if secret == "" || host == "" {
		return Conf{}, errors.New("Some configs are missing, trying reconfigure.")
	} else {
		return Conf{Secret : secret, Host: host,}, nil
	}
}

func getConfigFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	home := usr.HomeDir
	path := filepath.FromSlash(home + "/" + ConfigFileName)
	return path, nil
}

func configure(secret, host string) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if fileExists(configFile) {
		fmt.Println("rcp appears to have been configured already or it's a file name conflict, do you want to replace it (CANNOT BE UNDONEI BACKUP THE FILE!)? [y/N]")

		var action string 
		fmt.Scanln(&action)
		if action == "y" {
			//never trust users
			backupFile(configFile)
			err = os.Remove(configFile)

			if err != nil {
				return err
			}

		} else {
			fmt.Println("The file integrity remains")
			return nil
		}
	}

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	content := strings.ReplaceAll(ConfigFileContent, "<secret>", secret)
	content = strings.ReplaceAll(content, "<host>", host)

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func backupFile(path string) {
	 //Read all the contents of the  original file
    bytesRead, _ := ioutil.ReadFile(path)

    //Copy all the contents to the desitination file
    ioutil.WriteFile(path+".backup", bytesRead, 0755)
}

func fileExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func copyy(text string, c Conf) error {
	client := NewRcpRestClient(c)
	return client.DoPut(text)
}

func paste(c Conf) error{
	client := NewRcpRestClient(c)
	text, err := client.DoGet()

	if err != nil {
		return err
	}

	fmt.Println(text)
	return nil
}

func erase(c Conf) error {
	client := NewRcpRestClient(c)
	return client.DoDelete()
}

