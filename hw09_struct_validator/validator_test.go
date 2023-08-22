package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID         string `json:"id" validate:"len:36"`
		Name       string
		Age        int             `validate:"min:18|max:50"`
		Reputation int             `validate:"min:0|max:99"`
		Email      string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role       UserRole        `validate:"in:admin,stuff"`
		Phones     []string        `validate:"len:11"`
		meta       json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:         "1234567890",
				Age:        51,
				Reputation: -1,
				Email:      "biba.boba@yandex.ru",
				Role:       "stuf",
				Phones:     []string{"8920142288"},
			},
			fmt.Errorf("ID:%w\nAge:%w\nReputation:%w\nEmail:%w\nRole:%w", errLenNotMatch, errGreaterThanMax, errLessThanMin, errRegexNotMatch, errNotInSet), //nolint:lll
		},
		{
			User{
				ID:         "123456789012345678901234567890123456",
				Age:        26,
				Reputation: 30,
				Email:      "bibaboba@yandex.ru",
				Role:       "stuff",
				Phones:     []string{"89201422888"},
			},
			nil,
		},
		{
			App{
				Version: "123456",
			},
			fmt.Errorf("Version:%w", errLenNotMatch),
		},
		{
			App{
				Version: "12345",
			},
			nil,
		},
		{
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			Response{
				Code: 228,
				Body: "qwertyuiop[]",
			},
			fmt.Errorf("Code:%w", errNotInSet),
		},
		{
			Response{
				Code: 404,
				Body: "qwertyuiop[]",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err != nil {
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.Nil(t, tt.expectedErr)
			}
		})
	}
}
