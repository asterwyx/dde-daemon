/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/dbus"
	"sort"
	"strconv"
)

func (m *DockManager) allocEntryId() string {
	num := m.entryCount
	m.entryCount++
	return fmt.Sprintf("e%dT%x", num, getCurrentTimestamp())
}

func (m *DockManager) markAppLaunched(appInfo *AppInfo) {
	if appInfo == nil {
		return
	}
	id := appInfo.GetId()
	if id == "" {
		logger.Warning("markAppLaunched failed, appInfo %v no id", appInfo)
		return
	}
	go func() {
		if m.launcher == nil {
			return
		}
		logger.Infof("mark app %q launched", id)
		m.launcher.MarkLaunched(id)
		recordFrequency(id)
	}()
}

func (m *DockManager) attachOrDetachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	showOnDock := m.isWindowRegistered(win) && m.clientList.Contains(win) &&
		winInfo.canShowOnDock()
	logger.Debugf("win %v showOnDock? %v", win, showOnDock)
	entry := winInfo.entry
	if entry != nil {
		if !showOnDock {
			m.detachWindow(winInfo)
		}
	} else {

		if winInfo.entryInnerId == "" {
			winInfo.entryInnerId, winInfo.appInfo = m.identifyWindow(winInfo)
			m.markAppLaunched(winInfo.appInfo)
		} else {
			logger.Debugf("win %v identified", win)
		}

		if showOnDock {
			m.attachWindow(winInfo)
		}
	}
}

func (m *DockManager) initClientList() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	winSlice := windowSlice(clientList)
	sort.Sort(winSlice)
	m.clientList = winSlice
	for _, win := range winSlice {
		m.registerWindow(win)
	}
}

func (m *DockManager) initDockedApps() {
	dockedApps := uniqStrSlice(m.DockedApps.Get())
	for _, app := range dockedApps {
		m.appendDockedApp(app)
	}
	m.saveDockedApps()
}

func (m *DockManager) installAppEntry(e *AppEntry) {
	// install on session D-Bus
	err := dbus.InstallOnSession(e)
	if err != nil {
		logger.Warning("Install AppEntry to dbus failed:", err)
		return
	}

	entryObjPath := dbus.ObjectPath(entryDBusObjPathPrefix + e.Id)
	logger.Debugf("insertAndInstallAppEntry %v", entryObjPath)
	index := m.Entries.IndexOf(e)
	if index >= 0 {
		dbus.Emit(m, "EntryAdded", entryObjPath, int32(index))
	}
}

func (m *DockManager) addAppEntry(entryInnerId string, appInfo *AppInfo, index int) (*AppEntry, bool) {
	logger.Debug("addAppEntry innerId:", entryInnerId)

	var entry *AppEntry
	isNewAdded := false
	if e := m.Entries.GetFirstByInnerId(entryInnerId); e != nil {
		logger.Debug("entry existed")
		entry = e
		if appInfo != nil {
			appInfo.Destroy()
		}
	} else {
		// cache desktop hash => desktop file path
		if appInfo != nil {
			m.desktopHashFileMapCacheManager.SetKeyValue(appInfo.innerId, appInfo.GetFilePath())
			m.desktopHashFileMapCacheManager.AutoSave()
		}
		logger.Debug("entry not existed, newAppEntry")
		entry = newAppEntry(m, entryInnerId, appInfo)
		m.Entries = m.Entries.Insert(entry, index)
		logger.Debugf("insert entry %v at %v", entry.Id, index)
		isNewAdded = true
	}
	return entry, isNewAdded
}

func (m *DockManager) appendDockedApp(app string) {
	logger.Infof("appendDockedApp %q", app)
	appInfo := NewDockedAppInfo(app)
	if appInfo == nil {
		logger.Warning("appendDockedApp failed: appInfo is nil")
		return
	}
	entry, isNewAdded := m.addAppEntry(appInfo.innerId, appInfo, -1)
	entry.setIsDocked(true)
	entry.updateMenu()
	if isNewAdded {
		entry.updateName()
		entry.updateIcon()
		m.installAppEntry(entry)
	}
}

func (m *DockManager) removeAppEntry(e *AppEntry) {
	for _, entry := range m.Entries {
		if entry == e {
			dbus.UnInstallObject(e)

			entryId := entry.Id
			logger.Info("removeAppEntry id:", entryId)
			m.Entries = m.Entries.Remove(e)
			e.destroy()
			dbus.Emit(m, "EntryRemoved", entryId)
			return
		}
	}
	logger.Warning("removeAppEntry failed, entry not found")
}

func (m *DockManager) identifyWindow(winInfo *WindowInfo) (string, *AppInfo) {
	logger.Debugf("identifyWindow: window id: %v, window hash %v", winInfo.window, winInfo.innerId)
	if winInfo.innerId == "" {
		logger.Debug("identifyWindow failed winInfo no innerId")
		return "", nil
	}
	desktopHash := m.desktopWindowsMapCacheManager.GetKeyByValue(winInfo.innerId)
	logger.Debug("identifyWindow: get desktop hash:", desktopHash)
	var appInfo *AppInfo
	if desktopHash != "" {
		appInfo = m.desktopHashFileMapCacheManager.GetAppInfo(desktopHash)
		logger.Debug("identifyWindow: get AppInfo by desktop hash:", appInfo)
	}

	if appInfo == nil {
		// cache fail
		if desktopHash != "" {
			logger.Warning("winHash->DesktopHash success, but DesktopHash->appInfo fail")
			m.desktopHashFileMapCacheManager.DeleteKey(desktopHash)
			m.desktopWindowsMapCacheManager.DeleteKeyValue(desktopHash, winInfo.innerId)
		}

		var canCache bool
		appInfo, canCache = m.getAppInfoFromWindow(winInfo)
		logger.Debug("identifyWindow: getAppInfoFromWindow:", appInfo)
		if appInfo != nil && canCache {
			m.desktopWindowsMapCacheManager.AddKeyValue(appInfo.innerId, winInfo.innerId)
			m.desktopHashFileMapCacheManager.SetKeyValue(appInfo.innerId, appInfo.GetFilePath())
		}
	}

	var entryInnerId string
	if appInfo != nil {
		entryInnerId = appInfo.innerId
		logger.Debug("Set entryInnerId to desktop hash")
	} else {
		entryInnerId = winInfo.innerId
		logger.Debug("Set entryInnerId to window hash")
	}

	m.desktopWindowsMapCacheManager.AutoSave()
	m.desktopHashFileMapCacheManager.AutoSave()
	return entryInnerId, appInfo
}

func (m *DockManager) attachWindow(winInfo *WindowInfo) {
	var appInfoCopy *AppInfo
	if winInfo.appInfo != nil {
		appInfoCopy = NewAppInfoFromFile(winInfo.appInfo.GetFilePath())
	}
	entry, isNewAdded := m.addAppEntry(winInfo.entryInnerId, appInfoCopy, -1)
	entry.windowMutex.Lock()
	defer entry.windowMutex.Unlock()

	entry.attachWindow(winInfo)
	entry.updateMenu()
	if isNewAdded {
		entry.updateName()
		entry.updateIcon()
		m.installAppEntry(entry)
	}
}

func (m *DockManager) detachWindow(winInfo *WindowInfo) {
	entry := winInfo.entry
	if entry == nil {
		return
	}
	winInfo.entry = nil
	entry.windowMutex.Lock()
	defer entry.windowMutex.Unlock()

	detached := entry.detachWindow(winInfo)
	if !detached {
		return
	}
	if !entry.hasWindow() && !entry.IsDocked {
		m.removeAppEntry(entry)
		return
	}
	entry.updateWindowTitles()
	entry.updateIcon()
	entry.updateMenu()
	entry.updateIsActive()
}

func (m *DockManager) getAppInfoFromWindow(winInfo *WindowInfo) (*AppInfo, bool) {
	win := winInfo.window
	var ai *AppInfo

	gtkAppId := winInfo.gtkAppId
	logger.Debug("Try gtkAppId", gtkAppId)
	if gtkAppId != "" {
		ai = NewAppInfo(gtkAppId)
		if ai != nil {
			logger.Debugf("Get AppInfo success gtk app id: %q", gtkAppId)
			return ai, true
		}
	}

	// env GIO_LAUNCHED_DESKTOP_FILE
	var launchedDesktopFile string
	logger.Debug("Try process env")
	if winInfo.process != nil {
		envVars, err := getProcessEnvVars(winInfo.process.pid)
		if err == nil {
			launchedDesktopFile = envVars["GIO_LAUNCHED_DESKTOP_FILE"]
			pidStr := envVars["GIO_LAUNCHED_DESKTOP_FILE_PID"]
			launchedDesktopFilePid, _ := strconv.ParseUint(pidStr, 10, 32)
			logger.Debugf("launchedDesktopFile: %q, pid: %v", launchedDesktopFile, launchedDesktopFilePid)
			if winInfo.process.pid != 0 &&
				uint(launchedDesktopFilePid) == winInfo.process.pid {
				ai = NewAppInfoFromFile(launchedDesktopFile)
				if ai != nil {
					logger.Debugf("Get AppInfo success pid equal launchedDesktopFile: %q", launchedDesktopFile)
					return ai, true
				}
			}
		}
	}

	// bamf
	desktop := getDesktopFromWindowByBamf(win)
	logger.Debug("Try bamf")
	if desktop != "" {
		ai = NewAppInfoFromFile(desktop)
		if ai != nil {
			logger.Debugf("Get AppInfo success bamf desktop: %q", desktop)
			return ai, false
		}
	}

	// 通常不由 desktop 文件启动的应用 bamf 识别容易失败
	winGuessAppId := winInfo.guessAppId(m.appIdFilterGroup)
	logger.Debug("Try filter group", winGuessAppId)
	if winGuessAppId != "" {
		ai = NewAppInfo(winGuessAppId)
		if ai != nil {
			logger.Debugf("Get AppInfo success winGuessAppId: %q", winGuessAppId)
			return ai, false
		}
	}

	// wmclass
	if winInfo.wmClass != nil {
		logger.Debug("Try wmclass instance")
		instance := winInfo.wmClass.Instance
		if instance != "" {
			ai = NewAppInfo(instance)
			if ai != nil {
				logger.Debugf("Get AppInfo success wmClass instance %q", instance)
				return ai, true
			}
		}

		logger.Debug("Try wmclass class")
		class := winInfo.wmClass.Class
		if class != "" {
			ai = NewAppInfo(class)
			if ai != nil {
				logger.Debugf("Get AppInfo success wmClass class %q", class)
				return ai, true
			}
		}
	}

	logger.Debug("Try env var launchedDesktopFile")
	if launchedDesktopFile != "" {
		ai = NewAppInfoFromFile(launchedDesktopFile)
		if ai != nil {
			logger.Debugf("Get AppInfo success launchedDesktopFile %q", launchedDesktopFile)
			return ai, false
		}
	}

	logger.Debug("Get AppInfo failed")
	return nil, false
}