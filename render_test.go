package curl2go

import (
	"testing"
)

func TestRender_Render(t *testing.T) {
	testCases := []struct {
		input *Relevant
	}{
		{
			input: &Relevant{
				URL:    "https://example.com",
				Method: "GET",
			},
		},
	}

	render := NewRender()
	for _, tc := range testCases {
		_, err := render.Render(tc.input)
		if err != nil {
			t.Fatal(err)
		}
	}
}
