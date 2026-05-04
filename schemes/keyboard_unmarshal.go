package schemes

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON decodes Keyboard.Buttons by switching on each button's
// `type` field into the appropriate concrete struct. Without this method
// neither encoding/json nor jsoniter can decode `[][]ButtonInterface`
// from an inbound `inline_keyboard` attachment, which crashes the entire
// update with `decode non empty interface: can not unmarshal into nil`.
func (k *Keyboard) UnmarshalJSON(data []byte) error {
	var raw struct {
		Buttons [][]json.RawMessage `json:"buttons"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("keyboard: unmarshal envelope: %w", err)
	}

	k.Buttons = make([][]ButtonInterface, len(raw.Buttons))
	for i, row := range raw.Buttons {
		decoded := make([]ButtonInterface, 0, len(row))

		for _, rawBtn := range row {
			btn, err := decodeButton(rawBtn)
			if err != nil {
				continue
			}

			decoded = append(decoded, btn)
		}

		k.Buttons[i] = decoded
	}

	return nil
}

func decodeButton(data []byte) (ButtonInterface, error) {
	var base Button
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("button base: %w", err)
	}

	switch base.Type {
	case LINK:
		var b LinkButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("link button: %w", err)
		}

		return b, nil
	case CALLBACK:
		var b CallbackButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("callback button: %w", err)
		}

		return b, nil
	case OPEN_APP:
		var b OpenAppButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("open-app button: %w", err)
		}

		return b, nil
	case MESSAGE:
		var b MessageButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("message button: %w", err)
		}

		return b, nil
	case CLIPBOARD:
		var b ClipboardButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("clipboard button: %w", err)
		}

		return b, nil
	case CONTACT:
		var b RequestContactButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("contact button: %w", err)
		}

		return b, nil
	case GEOLOCATION:
		var b RequestGeoLocationButton
		if err := json.Unmarshal(data, &b); err != nil {
			return nil, fmt.Errorf("geo button: %w", err)
		}

		return b, nil
	}

	return base, nil
}