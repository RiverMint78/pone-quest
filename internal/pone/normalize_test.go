package pone

import "testing"

func TestNormalizeSearchText(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "preserve case and symbols",
			in:   "Hello@World #42!",
			want: "Hello@World #42!",
		},
		{
			name: "keep punctuation as-is",
			in:   "A.B, C? D! E-F",
			want: "A.B, C? D! E-F",
		},
		{
			name: "normalize quotes and collapse spaces",
			in:   "  “Twilight”   said:  'Hi'  ",
			want: "\"Twilight\" said: 'Hi'",
		},
		{
			name: "convert fullwidth to halfwidth",
			in:   "Ｈｅｌｌｏ，　ｗｏｒｌｄ！ １２３",
			want: "Hello, world! 123",
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
