// Code generated by "dbusutil-gen em -type ScreenSaver"; DO NOT EDIT.

package screensaver

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

func (v *ScreenSaver) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "Inhibit",
			Fn:      v.Inhibit,
			InArgs:  []string{"name", "reason"},
			OutArgs: []string{"cookie"},
		},
		{
			Name:   "SetTimeout",
			Fn:     v.SetTimeout,
			InArgs: []string{"seconds", "interval", "blank"},
		},
		{
			Name: "SimulateUserActivity",
			Fn:   v.SimulateUserActivity,
		},
		{
			Name:   "UnInhibit",
			Fn:     v.UnInhibit,
			InArgs: []string{"cookie"},
		},
	}
}
