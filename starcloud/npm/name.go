/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package npm

import (
	"fmt"
	"strings"
)

type NPMPackageIdentifier struct {
	Name string // Name of the package, may be of type myname or @org/myname
	Version string // Version number (the bit after the last @), e.g. 0.1.3
}

func PathToNPMPackage(path string) (NPMPackageIdentifier, error) {
	parts := strings.Split(path, "/")

	if len(parts) < 1 { // Should be at least package@1.2.3/, so two parts
		return NPMPackageIdentifier{}, fmt.Errorf("invalid npm filepath")
	}

	// NPM package with @org prefix
	if len(parts) > 1 && strings.HasPrefix(parts[0], "@") {
		return PackageIdToPackageAndVersion(parts[0] + "/" + parts[1])
	} else {
		return PackageIdToPackageAndVersion(parts[0])
	}

}

// starboard-notebook@1.2.3 -> starboard-notebook 1.2.3
func PackageIdToPackageAndVersion(pn string) (NPMPackageIdentifier, error) {
	parts := strings.Split(pn, "@")
	n := len(parts)

	if n < 2 {
		return NPMPackageIdentifier{}, fmt.Errorf("package %s doesn't have version in it", pn)
	}

	// At most two @ should be in a package id, e.g.
	// @org/package-name@1.2.3
	// package-name@1.2.3
	if n > 4 {
		return NPMPackageIdentifier{}, fmt.Errorf("package %s has more than 2 @ in its identifier", pn)
	}

	if n == 2 {
		return NPMPackageIdentifier{
			Name: parts[0],
			Version: parts[1],
		}, nil
	}

	// n = 3
	return NPMPackageIdentifier{
		Name: parts[0] + "@" + parts[1],
		Version: parts[2],
	}, nil
}