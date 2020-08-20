package worker

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func gitClone(dir string) {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: "https://github.com/tahirali-csc/hello-app",
		Auth: &http.BasicAuth{
			Username: "tahirali-csc",
			Password: "Tonyboy1:",
		},
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println(err)
	}
}
