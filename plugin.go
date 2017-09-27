package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type (
	Config struct {
		Hosts     []string `json:"host"`
		User      string   `json:"user"`
		Port      int      `json:"port"`
		Source    string   `json:"source"`
		Target    string   `json:"target"`
		Delete    bool     `json:"delete"`
		Recursive bool     `json:"recursive"`
		Include   []string `json:"include"`
		Exclude   []string `json:"exclude"`
		Filter    []string `json:"filter"`
		Commands  []string `json:"commands"`
		Key       string   `json:"key"`
	}

	Plugin struct {
		Config Config
	}
)

func (p Plugin) Exec() error {
	// write the rsa private key if provided
	if err := writeKey(p.Config.Key); err != nil {
		return err
	}

	// execute for each host
	for _, host := range p.Config.Hosts {
		// sync the files on the remote machine
		rs := p.buildRsync(host, "./") // TODO
		rs.Stderr = os.Stderr
		rs.Stdout = os.Stdout
		trace(rs)
		err := rs.Run()
		if err != nil {
			return err
		}

		// continue if no commands
		if len(p.Config.Commands) == 0 {
			continue
		}

		// execute commands on remote server (reboot instance, etc)
		if err := p.run(p.Config.Key, host); err != nil {
			return err
		}
	}

	return nil
}

// Build rsync command
func (p Plugin) buildRsync(host, root string) *exec.Cmd {

	var args []string
	args = append(args, "-az")

	// append recursive flag
	if p.Config.Recursive {
		args = append(args, "-r")
	}

	// append delete flag
	if p.Config.Delete {
		args = append(args, "--del")
	}

	// append custom ssh parameters
	args = append(args, "-e")
	args = append(args, fmt.Sprintf("'ssh -p %d -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -o StrictHostKeyChecking=no'", p.Config.Port))

	// append filtering rules
	for _, pattern := range p.Config.Include {
		args = append(args, fmt.Sprintf("--include=%s", pattern))
	}
	for _, pattern := range p.Config.Exclude {
		args = append(args, fmt.Sprintf("--exclude=%s", pattern))
	}
	for _, pattern := range p.Config.Filter {
		args = append(args, fmt.Sprintf("--filter=%s", pattern))
	}

	args = append(args, p.globSource(root)...)
	args = append(args, fmt.Sprintf("%s@%s:%s", p.Config.User, host, p.Config.Target))

	var shArgs []string
	var argBuf bytes.Buffer
	for _, pattern := range args {
		argBuf.WriteString(pattern)
		argBuf.WriteString(" ")
	}

	shArgs = append(shArgs, "-c")
	shArgs = append(shArgs, fmt.Sprintf("rsync %s", argBuf.String()))

	return exec.Command("/bin/sh", shArgs...)
}

// Run commands on the remote host
func (p Plugin) run(key, host string) error {

	// join the host and port if necessary
	addr := net.JoinHostPort(host, strconv.Itoa(p.Config.Port))

	// trace command used for debugging in the build logs
	fmt.Printf("$ ssh %s@%s -p %d \n", p.Config.User, addr, p.Config.Port)

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return fmt.Errorf("Error parsing private key. %s.", err)
	}

	config := &ssh.ClientConfig{
		User:            p.Config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("Error dialing server. %s.", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Error starting ssh session. %s.", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return session.Run(strings.Join(p.Config.Commands, "\n"))
}

// globSource returns the names of all files matching the source pattern.
// If there are no matches or an error occurs, the original source string is
// returned.
//
// If the source path is not absolute the root path will be prepended to the
// source path prior to matching.
func (p Plugin) globSource(root string) []string {
	src := p.Config.Source
	if !path.IsAbs(p.Config.Source) {
		src = path.Join(root, p.Config.Source)
	}
	srcs, err := filepath.Glob(src)
	if err != nil || len(srcs) == 0 {
		return []string{p.Config.Source}
	}
	sep := fmt.Sprintf("%c", os.PathSeparator)
	if strings.HasSuffix(p.Config.Source, sep) {
		// Add back the trailing slash removed by path.Join()
		for i := range srcs {
			srcs[i] += sep
		}
	}
	return srcs
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

// Writes the RSA private key
func writeKey(key string) error {
	if len(key) == 0 {
		return nil
	}

	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	sshpath := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return err
	}
	confpath := filepath.Join(sshpath, "config")
	privpath := filepath.Join(sshpath, "id_rsa")
	ioutil.WriteFile(confpath, []byte("StrictHostKeyChecking no\n"), 0700)
	return ioutil.WriteFile(privpath, []byte(key), 0600)
}
