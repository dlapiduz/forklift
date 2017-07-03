package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type RepoOptions struct {
	URL    string
	Branch string
	Token  string
}

func (ro *RepoOptions) CloneURL() string {
	return ro.URL
}

func main() {
	out, err := GitRebaseUpstream(
		&RepoOptions{
			URL:    "git@github.com:dlapiduz/deploy-to-cf.git",
			Branch: "new",
		},
		&RepoOptions{
			URL:    "https://github.com/jmcarp/deploy-to-cf.git",
			Branch: "master",
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)

}

func GitRebaseUpstream(repo, upstream *RepoOptions) (string, error) {

	dir, err := ioutil.TempDir("", "lift")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir) // clean up

	commands := []string{
		"git clone --depth 100 " + repo.CloneURL() + " .",
		"git checkout -b " + repo.Branch,
		"git remote add upstream " + upstream.CloneURL(),
		"git fetch upstream",
		"git rebase upstream/" + upstream.Branch,
		"git push -f origin " + repo.Branch,
	}

	var output string
	for _, c := range commands {
		args := strings.Split(c, " ")
		out, err := runCommand(dir, args[0], args[1:]...)
		if err != nil {
			return "", err
		}
		output = output + out + "\n"
	}

	return output, nil
}

func runCommand(dir string, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), err
	}

	return out.String(), nil
}
