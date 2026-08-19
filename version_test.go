package syngo_test

import (
	"github.com/stretchr/testify/require"
	"github.com/synesissoftware/syngo"

	"testing"
)

const (
	Expected_VersionMajor uint16 = 0
	Expected_VersionMinor uint16 = 1
	Expected_VersionPatch uint16 = 0
	Expected_VersionAB    uint16 = 0xFFFF
)

func Test_Version_Elements(t *testing.T) {
	require.Equal(t, Expected_VersionMajor, syngo.VersionMajor)
	require.Equal(t, Expected_VersionMinor, syngo.VersionMinor)
	require.Equal(t, Expected_VersionPatch, syngo.VersionPatch)
	require.Equal(t, Expected_VersionAB, syngo.VersionAB)
}

func Test_Version(t *testing.T) {
	require.Equal(t, uint64(0x0000_0001_0000_FFFF), syngo.Version)
}

func Test_Version_String(t *testing.T) {
	require.Equal(t, "0.1.0", syngo.VersionString())
}
