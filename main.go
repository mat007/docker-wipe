package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Microsoft/hcsshim"
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
	properties, err := hcsshim.GetContainers(hcsshim.ComputeSystemQuery{})
	if err != nil {
		return err
	}
	for _, property := range properties {
		if property.IsRuntimeTemplate {
			if err := remove(property, filepath.Join(dir, "containers")); err != nil {
				return err
			}
			if err := remove(property, filepath.Join(dir, "windowsfilter")); err != nil {
				return err
			}
		}
	}
	return nil
}

func remove(property hcsshim.ContainerProperties, dir string) error {
	dirs, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	for _, subdir := range dirs {
		if !strings.Contains(property.RuntimeImagePath, filepath.Join(dir, subdir.Name())) {
			continue
		}
		fmt.Println("Removing", property.ID)
		container, err := hcsshim.OpenContainer(property.ID)
		if err != nil {
			return err
		}
		if err := container.Terminate(); err != nil {
			return err
		}
		if err := container.Close(); err != nil {
			return err
		}
	}
	fmt.Println("Removing", dir)
	return os.RemoveAll(dir)
}
