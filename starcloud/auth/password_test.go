package auth

import "testing"

func TestHashAndVerifyPassword(t *testing.T) {
	tests := []struct {
		name            string
		password        string
		passwordAttempt string
	}{
		{
			"correct-password",
			"my-password",
			"my-password",
		},
		{
			"empty-password",
			"",
			"",
		},
		{
			"empty-password-nonempty-attempt",
			"",
			"bla",
		},
		{
			"empty-password-nonempty-attempt",
			"bla",
			"",
		},
		{
			"very-long-password",
			"asdasd1232341234324cvbhjyutrdfghdfghertywefsdfgsdfgsdfg34532sdfgfgert6eyrt sdfoj sfdgi10m93248120-3413ojisdlvdos!_%'ifjgowejrtpo123123vbfgfgd",
			"asdasd1232341234324cvbhjyutrdfghdfghertywefsdfgsdfgsdfg34532sdfgfgert6eyrt sdfoj sfdgi10m93248120-3413ojisdlvdos!_%'ifjgowejrtpo123123vbfgfgd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodedHash, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, false)
				return
			}

			valid, err := VerifyPassword(tt.passwordAttempt, encodedHash)
			if err != nil {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, false)
				return
			}

			if (tt.password == tt.passwordAttempt) != valid {
				t.Errorf("Verified different passwords, but they matched: %v %v", tt.password, tt.passwordAttempt)
			}
		})
	}
}
