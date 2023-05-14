package rdns

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLabelMap(t *testing.T) {
	m := make(labelTable)
	m.add([]string{"com", "example"})
	m.add([]string{"com"})
	require.Equal(t, 1, len(m))
	require.True(t, m["com"] != nil && len(m["com"]) == 0, "m[\"com\"] should be a zero-length map")

	m = make(labelTable)
	m.add([]string{"com"})
	m.add([]string{"com", "example"})
	require.Equal(t, 1, len(m))
	require.True(t, m["com"] != nil && len(m["com"]) == 0, "m[\"com\"] should be a zero-length map")

	m = make(labelTable)
	m.add([]string{"com", "example", "www"})
	m.add([]string{"nl", "example"})
	require.Len(t, m, 2)
	require.False(t, m.test([]string{"com", "example"}))
	require.True(t, m.test([]string{"com", "example", "www"}))
	require.True(t, m.test([]string{"com", "example", "www", "foo"}))
	require.False(t, m.test([]string{"www"}))
	require.False(t, m.test([]string{"www", "example", "com"}))
	require.True(t, m.test([]string{"nl", "example"}))
}
