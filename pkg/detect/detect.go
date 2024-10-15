package detect

import (
	"os"
	"sort"
	"strings"

	"github.com/zhanglongx/icon-detect/pkg/logfile"
	"golang.org/x/sys/windows/registry"
)

var (
	BOOST = []string{
		"Tortoise1Normal",
		"Tortoise2Modified",
		"Tortoise3Conflict",
		"Tortoise6Deleted",
		"Tortoise7Added",
		"Tortoise8Ignored",
		"Tortoise9Unversioned",
		"DropboxExt01",
		"DropboxExt02",
		"DropboxExt07",
		"OneDrive4",
	}

	KEY = `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\ShellIconOverlayIdentifiers`
)

type IconDetect struct {
	backup  map[string]string
	names   map[string]string
	deletes []string
	rename  map[string]string
}

func NewIconDetect() *IconDetect {
	return &IconDetect{
		backup:  make(map[string]string),
		names:   make(map[string]string),
		deletes: make([]string, 0),
		rename:  make(map[string]string),
	}
}

// Detect should be ran as a new IconDetect instance, and only once
func (i *IconDetect) Detect() (bool, error) {
	isChanged := false

	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		KEY,
		registry.READ)
	if err != nil {
		return isChanged, err
	}

	defer key.Close()

	names, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return isChanged, err
	}

	for _, n := range names {
		sub, err := registry.OpenKey(key, n, registry.READ)
		if err != nil {
			logfile.Printf("error opening key: %s, skip", n)
			continue
		}

		v, _, err := sub.GetStringValue("")

		if err != nil {
			logfile.Printf("error reading value: %s, skip", n)
			continue
		}

		sub.Close()

		i.backup[n] = v

		core := strings.TrimSpace(n)

		if _, ok := i.names[core]; ok {
			i.deletes = append(i.deletes, n)
			isChanged = true
		} else {
			i.names[core] = v
			if isIn(BOOST, core) {
				core = " " + core
			}
			if core != n {
				i.rename[n] = core
				isChanged = true
			}
		}
	}

	return isChanged, nil
}

// write backup into file with name as fileName, in .reg format
func (i *IconDetect) WriteBackup(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	file.WriteString("Windows Registry Editor Version 5.00\r\n\r\n")

	file.WriteString(`[HKEY_LOCAL_MACHINE\` + KEY + "]" + "\r\n\r\n")

	keys := make([]string, 0, len(i.backup))
	for k := range i.backup {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := i.backup[k]

		file.WriteString(`[HKEY_LOCAL_MACHINE\` + KEY + `\` + k + "]" + "\r\n")
		file.WriteString("@=\"" + v + "\"\r\n\r\n")
	}

	return nil
}

// Fix should be ran as a new IconDetect instance, and only once
func (i *IconDetect) Fix() error {
	for _, n := range i.deletes {
		err := registry.DeleteKey(registry.LOCAL_MACHINE,
			KEY+"\\"+n)

		if err != nil {
			logfile.Printf("error deleting key: %s, skip", n)
			continue
		}
	}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		KEY,
		registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer key.Close()

	for o, n := range i.rename {
		sub, err := registry.OpenKey(key, o, registry.ALL_ACCESS)
		if err != nil {
			logfile.Printf("error opening key: %s, skip", o)
			continue
		}

		v, _, err := sub.GetStringValue("")
		if err != nil {
			logfile.Printf("error reading value: %s, skip", o)
			continue
		}

		sub.Close()

		new, _, err := registry.CreateKey(key, n, registry.ALL_ACCESS)
		if err != nil {
			logfile.Printf("error creating key: %s, skip", n)
			continue
		}

		new.SetStringValue("", v)

		new.Close()

		logfile.Printf("rename \"%s\" to \"%s\"", o, n)
	}

	for o := range i.rename {
		err = registry.DeleteKey(registry.LOCAL_MACHINE, KEY+"\\"+o)
		if err != nil {
			logfile.Printf("error deleting key: %s, skip", o)
			continue
		}
	}

	return nil
}

func isIn(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
