// Code generated by "dbusutil-gen em -type Audio,Sink,SinkInput,Source,Meter"; DO NOT EDIT.

package audio

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

func (v *Audio) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "IsPortEnabled",
			Fn:      v.IsPortEnabled,
			InArgs:  []string{"cardId", "portName"},
			OutArgs: []string{"enabled"},
		},
		{
			Name: "NoRestartPulseAudio",
			Fn:   v.NoRestartPulseAudio,
		},
		{
			Name: "Reset",
			Fn:   v.Reset,
		},
		{
			Name:   "SetBluetoothAudioMode",
			Fn:     v.SetBluetoothAudioMode,
			InArgs: []string{"mode"},
		},
		{
			Name:   "SetPort",
			Fn:     v.SetPort,
			InArgs: []string{"cardId", "portName", "direction"},
		},
		{
			Name:   "SetPortEnabled",
			Fn:     v.SetPortEnabled,
			InArgs: []string{"cardId", "portName", "enabled"},
		},
	}
}
func (v *Meter) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name: "Tick",
			Fn:   v.Tick,
		},
	}
}
func (v *Sink) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "GetMeter",
			Fn:      v.GetMeter,
			OutArgs: []string{"meter"},
		},
		{
			Name:   "SetBalance",
			Fn:     v.SetBalance,
			InArgs: []string{"value", "isPlay"},
		},
		{
			Name:   "SetFade",
			Fn:     v.SetFade,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetMute",
			Fn:     v.SetMute,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetPort",
			Fn:     v.SetPort,
			InArgs: []string{"name"},
		},
		{
			Name:   "SetVolume",
			Fn:     v.SetVolume,
			InArgs: []string{"value", "isPlay"},
		},
	}
}
func (v *SinkInput) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:   "SetBalance",
			Fn:     v.SetBalance,
			InArgs: []string{"value", "isPlay"},
		},
		{
			Name:   "SetFade",
			Fn:     v.SetFade,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetMute",
			Fn:     v.SetMute,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetVolume",
			Fn:     v.SetVolume,
			InArgs: []string{"value", "isPlay"},
		},
	}
}
func (v *Source) GetExportedMethods() dbusutil.ExportedMethods {
	return dbusutil.ExportedMethods{
		{
			Name:    "GetMeter",
			Fn:      v.GetMeter,
			OutArgs: []string{"meter"},
		},
		{
			Name:   "SetBalance",
			Fn:     v.SetBalance,
			InArgs: []string{"value", "isPlay"},
		},
		{
			Name:   "SetFade",
			Fn:     v.SetFade,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetMute",
			Fn:     v.SetMute,
			InArgs: []string{"value"},
		},
		{
			Name:   "SetPort",
			Fn:     v.SetPort,
			InArgs: []string{"name"},
		},
		{
			Name:   "SetVolume",
			Fn:     v.SetVolume,
			InArgs: []string{"value", "isPlay"},
		},
	}
}
