// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2020/7/13

package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_attrAppNameFromOsEnv(t *testing.T) {
	_ = os.Unsetenv("CONTAINER_ID")
	_ = os.Unsetenv("APPSPACE_NAMESPACE")
	_ = os.Unsetenv("APPSPACE_IDC_NAME")

	t.Run("no env", func(t *testing.T) {
		t.Setenv("CONTAINER_ID", "")
		require.Empty(t, attrAppNameFromOsEnv())
	})

	t.Run("has env", func(t *testing.T) {
		t.Setenv("CONTAINER_ID", "1000004.bdapp-gdp-website-tucheng")
		require.Equal(t, "bdapp-gdp-website-tucheng", attrAppNameFromOsEnv())

		t.Setenv("CONTAINER_ID", "")
		t.Setenv("APPSPACE_NAMESPACE", "bdapp-gdp-website-1")
		t.Setenv("APPSPACE_IDC_NAME", "bjtest")
		require.Equal(t, "bdapp-gdp-website-1-bjtest", attrAppNameFromOsEnv())
	})

	t.Run("invalid env", func(t *testing.T) {
		t.Setenv("CONTAINER_ID", "bdapp-gdp-website-tucheng")
		require.Empty(t, attrAppNameFromOsEnv())

		t.Setenv("CONTAINER_ID", "1000004")
		require.Empty(t, attrAppNameFromOsEnv())
	})
}
