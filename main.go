package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	programData := os.Getenv("ProgramData")
	if err := wipe(filepath.Join(programData, "Docker")); err != nil {
		return err
	}
	dockerDesktop := filepath.Join(programData, "DockerDesktop", "vp-data-roots")
	roots, err := ioutil.ReadDir(dockerDesktop)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, root := range roots {
		if err := wipe(filepath.Join(dockerDesktop, root.Name())); err != nil {
			return err
		}
	}
	return nil
}

func wipe(dir string) error {
	fmt.Println("Wiping", dir)
	return nil
}
