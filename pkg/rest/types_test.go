package rest

import (
	"encoding/json"
	"testing"
)

func TestOptional(t *testing.T) {
	type T struct {
		Text   Optional[string]  `json:"text,omitempty"`
		Number Optional[float32] `json:"number,omitempty"`
		Bool   Optional[*bool]   `json:"bool,omitempty"`
	}

	tests := []struct {
		json           string
		text           string
		textDefined    bool
		number         float32
		numberDefined  bool
		boolean        *bool
		booleanDefined bool
	}{
		{
			json: `{}`,
		},
		{
			json:        `{"text": "Hi!"}`,
			textDefined: true,
			text:        "Hi!",
		},
		{
			json:          `{"text": "Hi!", "number": 23.4}`,
			textDefined:   true,
			text:          "Hi!",
			number:        23.4,
			numberDefined: true,
		},
		{
			json:           `{"bool": null}`,
			boolean:        nil,
			booleanDefined: true,
		},
	}

	for _, test := range tests {
		value := T{}
		err := json.Unmarshal([]byte(test.json), &value)

		if err != nil {
			t.Fatalf("Error parsing JSON: %v", err)
		}

		if test.textDefined != value.Text.Defined {
			t.Errorf("T.Text parsing error. Defined not matching: %v vs %v", test.text, value.Text.Value)
		} else if test.textDefined && test.text != value.Text.Value {
			t.Errorf("T.Text parsing error: %v vs %v", test.text, value.Text.Value)
		}

		if test.numberDefined != value.Number.Defined {
			t.Errorf("T.Number parsing error. Defined not matching: %v vs %v", test.number, value.Number.Value)
		} else if test.numberDefined && test.number != value.Number.Value {
			t.Errorf("T.Number parsing error: %v vs %v", test.number, value.Number.Value)
		}

		if test.booleanDefined != value.Bool.Defined {
			t.Errorf("T.Bool parsing error. Defined not matching: %v vs %v", test.boolean, value.Bool.Value)
		} else if test.booleanDefined && test.boolean != value.Bool.Value {
			t.Errorf("T.Bool parsing error: %v vs %v", test.boolean, value.Bool.Value)
		}

		// actual, err2 := json.Marshal(value)
		// if err2 != nil {
		// 	t.Fatalf("Error converting to JSON: %v", err2)
		// }

		// if string(actual) != test.json {
		// 	t.Errorf("Converted json mismatch: %v vs %v", string(actual), test.json)
		// }
	}
}
