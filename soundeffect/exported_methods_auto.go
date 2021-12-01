// Code generated by "dbusutil-gen em -type Manager"; DO NOT EDIT.

package soundeffect

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

func (v *Manager) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:   "EnableSound",
			Fn:     v.EnableSound,
			InArgs: []string{"name", "enabled"},
		},
		{
			Name:    "GetSoundEnabledMap",
			Fn:      v.GetSoundEnabledMap,
			OutArgs: []string{"result"},
		},
		{
			Name:    "GetSoundFile",
			Fn:      v.GetSoundFile,
			InArgs:  []string{"name"},
			OutArgs: []string{"file"},
		},
		{
			Name:    "GetSystemSoundFile",
			Fn:      v.GetSystemSoundFile,
			InArgs:  []string{"name"},
			OutArgs: []string{"file"},
		},
		{
			Name:    "IsSoundEnabled",
			Fn:      v.IsSoundEnabled,
			InArgs:  []string{"name"},
			OutArgs: []string{"enabled"},
		},
		{
			Name:   "PlaySound",
			Fn:     v.PlaySound,
			InArgs: []string{"name"},
		},
		{
			Name:   "PlaySystemSound",
			Fn:     v.PlaySystemSound,
			InArgs: []string{"name"},
		},
	}
}
