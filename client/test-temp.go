package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func createTempFile(contents string, jobFile string) string {
	jobDir, err := ioutil.TempDir("", "job")
	if err != nil {
		fmt.Println("Directory error", err)
		return ""
	}

	fullPath := path.Join(jobDir, jobFile)
	fmt.Println("Creating file::" + fullPath)

	f, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("File creation error", err)
		return ""
	}

	f.Write([]byte(contents))
	return jobDir
}

func main1() {
	dir := createTempFile("hello world\n", "hello.sh")
	fmt.Println(dir)

}
