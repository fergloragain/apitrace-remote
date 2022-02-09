package operations

import (
	"bytes"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func clonePublicRepo(repoURL, targetDirectory, branch string) (string, error) {

	var buf bytes.Buffer

	_, err := git.PlainClone(targetDirectory, false, &git.CloneOptions{
		URL:           repoURL,
		Progress:      &buf,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func clonePrivateRepo(auth ssh2.AuthMethod, repoURL, targetDirectory, branch string) (string, error) {
	var buf bytes.Buffer

	_, err := git.PlainClone(targetDirectory, false, &git.CloneOptions{
		URL:           repoURL,
		Progress:      &buf,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil

}
