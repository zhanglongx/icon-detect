package detect

import "gopkg.in/toast.v1"

func PushNotify(appID, title, scheme string) error {
	n := toast.Notification{
		AppID: appID,
		Title: title,
	}

	if scheme != "" {
		n.Actions = []toast.Action{{
			Type:      "protocol",
			Label:     "Restart",
			Arguments: scheme,
		}}
	}

	err := n.Push()
	if err != nil {
		return err
	}

	return nil
}
