package model

import (
	"testing"

	"github.com/misoboy/cdule/pkg"
	"github.com/stretchr/testify/require"
)

func Test_ReadConfigInMemory(t *testing.T) {
	param := []string{"../../resources", "config_in_memory"}
	cduleConfig, err := readConfig(param)
	require.NoError(t, err)
	require.Equal(t, string(pkg.MEMORY), cduleConfig.Cdule.Type)
	require.Equal(t, "/Users/dsinghvi/sqlite.db", cduleConfig.Database.Url)
}

func Test_ReadConfigInDB(t *testing.T) {
	param := []string{"../../resources", "config"}
	cduleConfig, err := readConfig(param)
	require.NoError(t, err)
	require.Equal(t, string(pkg.DATABASE), cduleConfig.Cdule.Type)
	require.Equal(t, "postgres://cduleuser:cdulepassword@localhost:5432/cdule?sslmode=disable", cduleConfig.Database.Url)
}

func Test_ReadConfig_InvalidConfigPath(t *testing.T) {
	param := []string{"../resourceswhichdoesnotexists", "invalidconfig"}
	_, err := readConfig(param)
	require.Error(t, err)
}
