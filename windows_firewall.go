// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func init() {
	register("WindowsFirewall", cleanWindowsFirewall)
}

func removeFromFwMgr(name, appname string) error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", fmt.Sprintf(`name="%s"`, name), fmt.Sprintf(`program="%s"`, appname))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cleanWindowsFirewall(dryrun, verbose bool) (string, error) {
	pattern := filepath.Join(os.TempDir(), "go-build")

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unk, err := oleutil.CreateObject("HNetCfg.FwPolicy2")
	if err != nil {
		return "", err
	}
	dsp, err := unk.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return "", err
	}
	rules := oleutil.MustGetProperty(dsp, "Rules").ToIDispatch()

	removed := 0
	err = oleutil.ForEach(rules, func(v *ole.VARIANT) error {
		rule := v.ToIDispatch()
		name := oleutil.MustGetProperty(rule, "Name").ToString()
		appname := oleutil.MustGetProperty(rule, "Applicationname").ToString()
		if strings.HasPrefix(strings.ToLower(appname), strings.ToLower(pattern)) {
			if verbose {
				log.Println("WindowsFirewall:", name)
			}
			if !dryrun {
				_, rerr := oleutil.CallMethod(rules, "Remove", name)
				if rerr != nil {
					rerr = removeFromFwMgr(name, appname)
					if rerr != nil {
						return rerr
					}
				}
			}
			removed++
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	result := "%dth rules removed"
	if dryrun {
		result = "%dth rules removable"
	}
	return fmt.Sprintf(result, removed), nil
}
