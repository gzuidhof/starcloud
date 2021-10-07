/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package account

import "testing"

func TestValidUsernames(t *testing.T) {
	tests := []string{
		"hello",
		"hello-there",
		"ab",
		"abc",
		"123",
		"h-e",
		"pretty-long-but-not-too-long",
		"AlphaNumeric4l-TeSt0",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if err := IsValidUsername(name); err != nil {
				t.Errorf("IsValidUsername() error = %v, valid true", err)
			}
		})
	}
}

func TestInvalidUsernames(t *testing.T) {
	tests := []string{
		"a", // Too short
		"z",
		"-",
		"--",
		"---",
		"a--",
		"aaa-",
		"aaa--",
		"-aaa",
		"aa_aa",
		"a?a",
		"aa--aa",
		"way-tooooooooooooooooooooooooooooooooo-long",

		// Restricted usernames
		"starboard",
		"admin",
		"host",
		"settings",
		"starboardsomething",
		"starboard",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			if err := IsValidUsername(name); err == nil {
				t.Errorf("IsValidUsername() no error, shouldn't have been valid")
			}
		})
	}
}
