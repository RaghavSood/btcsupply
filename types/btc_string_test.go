package types

import "testing"

func Test_BTCString_String(t *testing.T) {
	b := BTCString("100.0")

	want := "100.0"
	got := b.String()

	if got != want {
		t.Errorf("BTCString.String() = %v; want %v", got, want)
	}
}

func Test_BTCString_UnmarshalJSON(t *testing.T) {
	var b BTCString

	var tests = []struct {
		have []byte
		want string
	}{
		{[]byte(`"100.0"`), "100.0"},
		{[]byte(`"0.00000000"`), "0.00000000"},
		{[]byte(`100.0`), "100.0"},
	}

	for _, test := range tests {
		err := b.UnmarshalJSON(test.have)
		if err != nil {
			t.Errorf("BTCString.UnmarshalJSON() = %v; want nil", err)
		}

		if string(b) != test.want {
			t.Errorf("BTCString.UnmarshalJSON() = %v; want %v", string(b), test.want)
		}
	}
}

func Test_BTCString_MarshalJSON(t *testing.T) {
	var tests = []struct {
		have BTCString
		want string
	}{
		{BTCString("100.0"), `"100.0"`},
		{BTCString("0.00000000"), `"0.00000000"`},
		{BTCString("100.0"), `"100.0"`},
	}

	for _, test := range tests {
		got, err := test.have.MarshalJSON()
		if err != nil {
			t.Errorf("BTCString.MarshalJSON() = %v; want nil", err)
		}

		if string(got) != test.want {
			t.Errorf("BTCString.MarshalJSON() = %v; want %v", string(got), test.want)
		}
	}
}

func Test_BTCString_NonZero(t *testing.T) {
	var tests = []struct {
		have BTCString
		want bool
	}{
		{BTCString("100.0"), true},
		{BTCString("0.00000000"), false},
		{BTCString("0"), false},
		{BTCString(""), false},
	}

	for _, test := range tests {
		got := test.have.NonZero()
		if got != test.want {
			t.Errorf("BTCString.NonZero() = %v; want %v", got, test.want)
		}
	}
}
