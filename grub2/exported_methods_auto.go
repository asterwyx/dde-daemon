// Code generated by "dbusutil-gen em -type Grub2,Theme,EditAuth"; DO NOT EDIT.

package grub2

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

func (v *EditAuth) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:   "Disable",
			Fn:     v.Disable,
			InArgs: []string{"username"},
		},
		{
			Name:   "Enable",
			Fn:     v.Enable,
			InArgs: []string{"username", "password"},
		},
	}
}
func (v *Grub2) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "GetAvailableGfxmodes",
			Fn:      v.GetAvailableGfxmodes,
			OutArgs: []string{"gfxModes"},
		},
		{
			Name:    "GetSimpleEntryTitles",
			Fn:      v.GetSimpleEntryTitles,
			OutArgs: []string{"titles"},
		},
		{
			Name: "PrepareGfxmodeDetect",
			Fn:   v.PrepareGfxmodeDetect,
		},
		{
			Name: "Reset",
			Fn:   v.Reset,
		},
		{
			Name:   "SetDefaultEntry",
			Fn:     v.SetDefaultEntry,
			InArgs: []string{"entry"},
		},
		{
			Name:   "SetEnableTheme",
			Fn:     v.SetEnableTheme,
			InArgs: []string{"enabled"},
		},
		{
			Name:   "SetGfxmode",
			Fn:     v.SetGfxmode,
			InArgs: []string{"gfxmode"},
		},
		{
			Name:   "SetTimeout",
			Fn:     v.SetTimeout,
			InArgs: []string{"timeout"},
		},
	}
}
func (v *Theme) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "GetBackground",
			Fn:      v.GetBackground,
			OutArgs: []string{"background"},
		},
		{
			Name:   "SetBackgroundSourceFile",
			Fn:     v.SetBackgroundSourceFile,
			InArgs: []string{"filename"},
		},
	}
}
