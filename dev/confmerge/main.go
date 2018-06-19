package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sourcegraph/sourcegraph/pkg/conf"
	"github.com/sourcegraph/sourcegraph/schema"
)

// merge multiple config files into a single config file
func main() {
	var cfg *schema.SiteConfiguration
	for _, path := range os.Args[1:] {
		cfgBuf, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("can't read '%s': %s\n", path, err)
			continue
		}
		cfgData, err := conf.ParseConfigData(string(cfgBuf))
		if err != nil {
			fmt.Printf("can't parse '%s': %s\n", path, err)
			continue
		}
		cfg = conf.AppendConfig(cfg, cfgData)
	}
	j, err := json.Marshal(cfg)
	if err != nil {
		// we can't really do anything if fmt.Errorf is failing, but
		// go vet requires it be checked.
		_ = fmt.Errorf("fatal: %s", err)
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, j, "", "  ")
	if err != nil {
		_ = fmt.Errorf("fatal: %s", err)
		return
	}
	out.WriteTo(os.Stdout)
	fmt.Println()
}
