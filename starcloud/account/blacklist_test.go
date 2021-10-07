/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package account

import "testing"

func Test_isInBlacklist(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			"should-be-in-blacklist",
			"admin",
			true,
		},
		{
			"should-not-be-in-blacklist-1",
			"adminblabla",
			false,
		},
		{
			"should-not-be-in-blacklist-2",
			"qwerty",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInBlacklist(tt.username); got != tt.want {
				t.Errorf("isInBlacklist() = %v, want %v", got, tt.want)
			}
		})
	}
}
