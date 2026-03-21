package pone

import "testing"

func TestNormalizeSearchText(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want  string
	}{
		{
			name: "lowercase and remove symbols",
			in:   "Hello@World #42!",
			want: "helloworld 42!",
		},
		{
			name: "keep allowed punctuation",
			in:   "A.B, C? D! E-F",
			want: "a.b, c? d! e-f",
		},
		{
			name: "normalize quotes and collapse spaces",
			in:   "  “Twilight”   said:  'Hi'  ",
			want: "\"twilight\" said \"hi\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSearchText(tt.in)
			if got != tt.want {
				t.Fatalf("NormalizeSearchText() = %q, want %q", got, tt.want)
			}
		})
	}
}
