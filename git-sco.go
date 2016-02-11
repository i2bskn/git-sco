package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const delimiter = "/"

var inFeatureNamespace = flag.Bool("f", false, "Omit the feature namespace")

func IsInsideWorkTree() bool {
	err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run()
	if err != nil {
		return false
	}
	return true
}

func LocalBranches() ([]string, error) {
	branches := make([]string, 0, 0)
	stdout, err := exec.Command("git", "branch").Output()
	if err != nil {
		return branches, err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		branch := lineTrim(strings.Replace(line, "*", "", 1))
		branches = append(branches, branch)
	}
	return branches, err
}

func RemoteBranches() (map[string][]string, error) {
	branches := make(map[string][]string)
	stdout, err := exec.Command("git", "branch", "--remotes").Output()
	if err != nil {
		return branches, err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		remote, branch := splitRemoteName(line)
		_, exist := branches[remote]
		if !exist {
			branches[remote] = make([]string, 0, 0)
		}

		branches[remote] = append(branches[remote], branch)
	}

	return branches, err
}

func splitRemoteName(line string) (string, string) {
	elements := strings.Split(lineTrim(line), delimiter)
	return elements[0], strings.Join(elements[1:], delimiter)
}

func lineTrim(line string) string {
	return strings.Trim(strings.TrimRight(line, "\n"), " ")
}

func main() {
	flag.Parse()

	var err error

	if len(os.Args) != 2 {
		fmt.Printf("%s <branch>\n", os.Args[0])
		os.Exit(1)
	}

	if !IsInsideWorkTree() {
		fmt.Println("Not a git repository")
		os.Exit(1)
	}

	branchName := os.Args[1]
	if *inFeatureNamespace {
		branchName = strings.Join([]string{"feature", branchName}, delimiter)
	}

	var localBranches []string
	localBranches, err = LocalBranches()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	isExist := false
	for _, branch := range localBranches {
		if branch == branchName {
			isExist = true
			break
		}
	}

	if isExist {
		err = exec.Command("git", "checkout", branchName).Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	var remoteBranches map[string][]string
	remoteBranches, err = RemoteBranches()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var remoteName string
	for remote, branches := range remoteBranches {
		for _, branch := range branches {
			if branch == branchName {
				remoteName = remote
				break
			}
		}

		if len(remoteName) != 0 {
			break
		}
	}

	var checkout *exec.Cmd
	if len(remoteName) != 0 {
		remoteBranchName := strings.Join([]string{remoteName, branchName}, delimiter)
		checkout = exec.Command("git", "checkout", "-b", branchName, remoteBranchName)
	} else {
		checkout = exec.Command("git", "checkout", "-b", branchName)
	}

	err = checkout.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
