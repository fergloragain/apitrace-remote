package operations

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
)

func Clone(user, privateKeyPath, repoURL, targetDirectory, branch string) (string, string, error) {

	if len(privateKeyPath) == 0 || len(user) == 0 {
		stdout, err := clonePublicRepo(repoURL, targetDirectory, branch)

		if err != nil {
			return stdout, "", err
		}

		return stdout, "", nil
	} else {
		pem, err := ioutil.ReadFile(privateKeyPath)

		if err != nil {
			return "", "", err
		}

		signer, err := ssh.ParsePrivateKey(pem)

		if err != nil {
			return "", "", err
		}

		auth := &ssh2.PublicKeys{User: user, Signer: signer}
		stdout, err := clonePrivateRepo(auth, repoURL, targetDirectory, branch)

		if err != nil {
			return stdout, "", err
		}

		return stdout, "", nil
	}
}

func Build(workingDirectory, buildScript string) (string, string, error) {

	args := []string{
		buildScript,
	}

	stdout, stderr, err := execute(workingDirectory, "/bin/sh", args)

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}

func Trace(workingDirectory, apiTraceLocation, executableToTrace string, timeout int) (string, string, error) {

	args := []string{
		fmt.Sprintf("%ds", timeout),
		apiTraceLocation,
		"trace",
		fmt.Sprintf("./%s", executableToTrace),
	}

	stdout, stderr, err := execute(workingDirectory, "timeout", args)

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}

func Dump(workingDirectory, apitraceLocation, traceLocation string) (string, string, error) {

	args := []string{
		"dump",
		"-v",
		"--color=never",
		traceLocation,
	}

	stdout, stderr, err := execute(workingDirectory, apitraceLocation, args)

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}

func DumpImages(workingDirectory, apitraceLocation, traceLocation, callID string) (string, string, error) {

	args := []string{
		"dump-images",
		"-m",
		fmt.Sprintf("--calls=%s", callID),
		traceLocation,
	}

	stdout, stderr, err := execute(workingDirectory, apitraceLocation, args)

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}

func Retrace(workingDirectory, glretraceLocation, traceLocation, callID string) (string, string, error) {

	args := []string{
		"-v",
		"--sb",
		fmt.Sprintf("--D=%s", callID),
		"--dump-format=json",
		traceLocation,
		"|",
		"tr",
		"-d",
		"'\n'",
	}

	stdout, _, err := execute(workingDirectory, glretraceLocation, args)

	if err != nil {
		return stdout, "", err
	}

	args = []string{
		"-v",
		"--sb",
		fmt.Sprintf("--D=%s", callID),
		"--dump-format=ubjson",
		traceLocation,
	}

	_, stderr, err := execute(workingDirectory, glretraceLocation, args)

	if err != nil {
		return stdout, stderr, err
	}

	return stdout, stderr, nil
}
