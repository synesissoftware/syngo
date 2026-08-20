package syngo_test

import (
	"github.com/synesissoftware/syngo"

	"github.com/stretchr/testify/require"

	"testing"
)

const (
	Expected_VersionMajor uint16 = 0
	Expected_VersionMinor uint16 = 3
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
	require.Equal(t, uint64(0x0000_0003_0000_FFFF), syngo.Version())
}

func Test_Version_String(t *testing.T) {
	require.Equal(t, "0.3.0", syngo.VersionString())
}
