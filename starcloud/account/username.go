/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package account

import (
	"fmt"
	"regexp"
	"strings"
)

var validCharactersRegex = regexp.MustCompile(`^[A-Za-z\d-]+$`)

func IsValidUsername(username string) error {
	if isInBlacklist(username) || strings.HasPrefix(username, "starb") || strings.HasPrefix(username, "starc") {
		return fmt.Errorf("username is restricted")
	} else if len(username) < 2 {
		return fmt.Errorf("username is too short")
	} else if len(username) > 39 {
		return fmt.Errorf("username is too long")
	} else if strings.HasPrefix(username, "-") {
		return fmt.Errorf("username can't start with a hyphen")
	} else if strings.HasSuffix(username, "-") {
		return fmt.Errorf("username can't end with a hyphen")
	} else if !validCharactersRegex.MatchString(username) {
		return fmt.Errorf("username contains invalid characters")
	} else if strings.Contains(username, "--") {
		return fmt.Errorf("username contains consecutive hyphens")
	}
	return nil
}
