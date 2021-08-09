/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package config

import (
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("port", 8080)
}
