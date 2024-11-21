package switchboardconfig

import (
	"github.com/NorskHelsenett/ror/pkg/config/configconsts"
	"github.com/NorskHelsenett/ror/pkg/config/rorversion"

	"github.com/NorskHelsenett/ror/pkg/clients/vaultclient"
	"github.com/NorskHelsenett/ror/pkg/rlog"

	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var (
	VaultClient *vaultclient.VaultClient
	VaultSecret *vault.Secret
	Version     string = "1.1.0"
	Commit      string = "FFFFF"
)

func init() {
	viper.AutomaticEnv()
	viper.SetDefault(configconsts.VERSION, Version)
	viper.SetDefault(configconsts.COMMIT, Commit)
	viper.SetDefault(configconsts.ENVIRONMENT, "production")
	viper.SetDefault(rlog.LOG_LEVEL, "info")
	viper.SetDefault(configconsts.ROLE, "ror-ms-switchboard")
}

func Load() {
	environment := viper.GetString(configconsts.ENVIRONMENT)
	rlog.Info("loaded environment", rlog.String("Environment", environment))

	_ = viper.WriteConfig()
}

func GetRorVersion() rorversion.RorVersion {
	return rorversion.NewRorVersion(viper.GetString(configconsts.VERSION), viper.GetString(configconsts.COMMIT))
}
