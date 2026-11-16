package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const instance = "shrimp"

var baseFlags = []string{"--instance", instance}

var files = []string{
	"mediamtx.yml",
	"static-dhcp-lease.json",
	"wifi.json",
}

func main() {
	// create base config if not exists
	_, err := os.Stat(confpath())
	if os.IsNotExist(err) {
		run("new")
		copy(confpath(), "gokrazy/config_base.json")
	} else if err != nil {
		panic(err)
	}

	// merge the base config with the public config
	base := readJSON("gokrazy/config_base.json")
	custom := readJSON("gokrazy/customizations.json")
	merged := deepMerge(base, custom)
	writeJSON(merged, confpath())

	for _, f := range files {
		copy(fmt.Sprintf("gokrazy/%v", f), fmt.Sprintf("%v/%v", confdir(), f))
	}

	if len(os.Args) < 2 {
		return // nothing to do
	}

	cmd := os.Args[1]
	switch cmd {
	case "overwrite":
		run(os.Args[1:]...)
	default:
		fmt.Printf("unknown command: %v\n", cmd)
		os.Exit(1)
	}
}

func run(cmd ...string) {
	var allargs []string
	allargs = append(allargs, baseFlags...)
	allargs = append(allargs, cmd...)
	fmt.Printf("gok %v\n", strings.Join(allargs, " "))
	c := exec.Command("gok", allargs...)
	out, err := c.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		panic(err)
	}
}

func copy(srcpath, destpath string) {
	fmt.Printf("cp %v %v\n", srcpath, destpath)
	data, err := os.ReadFile(srcpath)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(destpath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func readJSON(path string) map[string]any {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var out map[string]any
	err = json.Unmarshal(data, &out)
	if err != nil {
		panic(err)
	}
	return out
}

func writeJSON(data map[string]any, path string) {
	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, j, 0644)
	if err != nil {
		panic(err)
	}
}

func deepMerge(dst, src map[string]any) map[string]any {
	for key, srcVal := range src {
		if dstVal, exists := dst[key]; exists {
			// Check if both values are maps (nested JSON objects)
			dstMap, dstIsMap := dstVal.(map[string]any)
			srcMap, srcIsMap := srcVal.(map[string]any)

			if dstIsMap && srcIsMap {
				// Recursively merge nested maps
				dst[key] = deepMerge(dstMap, srcMap)
			} else {
				// Otherwise, overwrite with the source value
				dst[key] = srcVal
			}
		} else {
			// Key doesn't exist in dst, add it
			dst[key] = srcVal
		}
	}
	return dst
}

func confpath() string {
	return fmt.Sprintf("%v/config.json", confdir())
}

func confdir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v/gokrazy/%v", home, instance)
}
