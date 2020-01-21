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
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Println(filepath.Base(os.Args[0]), "<path>")
		os.Exit(1)
	}
	if err := wipe(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var driverInfo = hcsshim.DriverInfo{}

func wipe(dir string) error {
	fmt.Println("Wiping", dir)
	properties, err := hcsshim.GetContainers(hcsshim.ComputeSystemQuery{})
	if err != nil {
		return err
	}
	if err := remove(properties, filepath.Join(dir, "containers")); err != nil {
		return err
	}
	if err := remove(properties, filepath.Join(dir, "windowsfilter")); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(dir, "image"))
}

func remove(properties []hcsshim.ContainerProperties, dir string) error {
	layers, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		fmt.Println("Skipping", dir, "(not found)")
		return nil
	}
	fmt.Println("Checking", dir)
	if err != nil {
		return err
	}
	for _, layer := range layers {
		path := filepath.Join(dir, layer.Name())
		for _, property := range properties {
			if !property.IsRuntimeTemplate || !strings.Contains(property.RuntimeImagePath, path) {
				continue
			}
			fmt.Println("Removing", path)
			container, err := hcsshim.OpenContainer(property.ID)
			if err != nil {
				return err
			}
			// Ignoring error as it can be asynchronous behind the scene.
			_ = container.Terminate()
			if err := container.Close(); err != nil {
				return err
			}
		}
		if err := hcsshim.DestroyLayer(driverInfo, path); err != nil {
			return err
		}
	}
	fmt.Println("Deleting", dir)
	return os.RemoveAll(dir)
}
