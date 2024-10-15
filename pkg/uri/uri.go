package uri

import (
	"golang.org/x/sys/windows/registry"
)

func IsURISchemeRegistered(scheme string) bool {
	key, err := registry.OpenKey(registry.CLASSES_ROOT, scheme, registry.QUERY_VALUE)
	if err != nil {
		return false
	}

	defer key.Close()

	return true
}

func RegisterURIScheme(app, scheme, path string) error {
	key, _, err := registry.CreateKey(registry.CLASSES_ROOT, scheme, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	if err := key.SetStringValue("", "URL:"+app+" Protocol"); err != nil {
		return err
	}

	if err := key.SetStringValue("URL Protocol", ""); err != nil {
		return err
	}

	cmdKey, _, err := registry.CreateKey(key, `shell\open\command`, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer cmdKey.Close()

	if err := cmdKey.SetStringValue("", `"`+path+`" "%1"`); err != nil {
		return err
	}

	return nil
}

func UnRegisterRIScheme(scheme string) error {
	key, err := registry.OpenKey(registry.CLASSES_ROOT, scheme,
		registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS|registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}

	defer key.Close()

	err = deleteSubKeys(key, scheme)
	if err != nil {
		return err
	}

	err = registry.DeleteKey(registry.CLASSES_ROOT, scheme)
	if err != nil {
		return err
	}

	return nil
}

func deleteSubKeys(key registry.Key, path string) error {
	subKeyNames, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return err
	}

	for _, subKeyName := range subKeyNames {
		subKeyPath := path + `\` + subKeyName
		subKey, err := registry.OpenKey(key, subKeyName, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			return err
		}

		err = deleteSubKeys(subKey, subKeyPath)
		subKey.Close()
		if err != nil {
			return err
		}

		err = registry.DeleteKey(key, subKeyName)
		if err != nil {
			return err
		}
	}

	return nil
}
