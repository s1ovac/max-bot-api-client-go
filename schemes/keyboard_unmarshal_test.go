package schemes

import (
	"encoding/json"
	"testing"
)

func TestKeyboardUnmarshal_LinkButton(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"buttons":[[{"url":"https://vk.com","text":"кнопка","type":"link"}]]}`)

	var kb Keyboard
	if err := json.Unmarshal(payload, &kb); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(kb.Buttons) != 1 || len(kb.Buttons[0]) != 1 {
		t.Fatalf("unexpected shape: %+v", kb.Buttons)
	}

	link, ok := kb.Buttons[0][0].(LinkButton)
	if !ok {
		t.Fatalf("expected LinkButton, got %T", kb.Buttons[0][0])
	}

	if link.Url != "https://vk.com" {
		t.Errorf("url = %q, want %q", link.Url, "https://vk.com")
	}

	if link.Text != "кнопка" {
		t.Errorf("text = %q, want %q", link.Text, "кнопка")
	}

	if link.Type != LINK {
		t.Errorf("type = %q, want %q", link.Type, LINK)
	}
}

func TestKeyboardUnmarshal_MultipleRowsAndTypes(t *testing.T) {
	t.Parallel()

	payload := []byte(`{
		"buttons":[
			[{"url":"https://vk.com","text":"link","type":"link"}],
			[{"payload":"cb1","text":"cb","type":"callback","intent":"positive"}]
		]
	}`)

	var kb Keyboard
	if err := json.Unmarshal(payload, &kb); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(kb.Buttons) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(kb.Buttons))
	}

	if _, ok := kb.Buttons[0][0].(LinkButton); !ok {
		t.Errorf("row 0: expected LinkButton, got %T", kb.Buttons[0][0])
	}

	if _, ok := kb.Buttons[1][0].(CallbackButton); !ok {
		t.Errorf("row 1: expected CallbackButton, got %T", kb.Buttons[1][0])
	}
}

func TestKeyboardUnmarshal_UnknownTypeFallsBackToBase(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"buttons":[[{"text":"unknown","type":"some_future_type"}]]}`)

	var kb Keyboard
	if err := json.Unmarshal(payload, &kb); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(kb.Buttons) != 1 || len(kb.Buttons[0]) != 1 {
		t.Fatalf("unexpected shape: %+v", kb.Buttons)
	}

	base, ok := kb.Buttons[0][0].(Button)
	if !ok {
		t.Fatalf("expected base Button fallback, got %T", kb.Buttons[0][0])
	}

	if base.Text != "unknown" {
		t.Errorf("text = %q, want %q", base.Text, "unknown")
	}
}

func TestKeyboardUnmarshal_EmptyButtons(t *testing.T) {
	t.Parallel()

	var kb Keyboard
	if err := json.Unmarshal([]byte(`{"buttons":[]}`), &kb); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(kb.Buttons) != 0 {
		t.Errorf("expected empty buttons, got %d rows", len(kb.Buttons))
	}
}