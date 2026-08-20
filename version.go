// Copyright 2019-2026 Matthew Wilson and Synesis Information Systems. All
// rights reserved. Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Created: 14th March 2019
 * Updated: 20th August 2026
 */

package syngo

import "github.com/synesissoftware/ver2go"

const (
	VersionMajor uint16 = 0
	VersionMinor uint16 = 3
	VersionPatch uint16 = 1
	VersionAB    uint16 = ver2go.Release
)

var (
	version              = ver2go.CombineVersion(VersionMajor, VersionMinor, VersionPatch, VersionAB)
	versionString string = ver2go.CalcVersionString(VersionMajor, VersionMinor, VersionPatch, VersionAB)
)

// Version returns this library's version as a packed 64-bit integer, formed
// by ver2go.CombineVersion from VersionMajor, VersionMinor, VersionPatch,
// and VersionAB. The result is suitable for numeric comparison: a later
// release has a strictly greater value than an earlier one that uses the
// same packing.
func Version() uint64 {
	return version
}

// VersionString returns this library's version as a human-readable string,
// formed by ver2go.CalcVersionString from VersionMajor, VersionMinor,
// VersionPatch, and VersionAB. For a final (non-prerelease) version the
// result is of the form "MAJOR.MINOR.PATCH", e.g. "0.3.1".
func VersionString() string {
	return versionString
}

/* ///////////////////////////// end of file //////////////////////////// */
