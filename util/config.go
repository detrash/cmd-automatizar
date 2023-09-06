package util

import "github.com/spf13/viper"

type Config struct {
	DBPass string `mapstructure:"DB_PASS"`
	DBUser string `mapstructure:"DB_USER"`
	DBHost string `mapstructure:"DB_HOST"`
	DBNAME string `mapstructure:"DB_NAME"`

	ProviderCelo string `mapstructure:"PROVIDER_CELO"`
	SecretCelo   string `mapstructure:"SECRET_CELO"`
	ProviderPol  string `mapstructure:"PROVIDER_POL"`

	SecretPol string `mapstructure:"SECRET_POL"`

	CRecyAddress string `mapstructure:"CRECY_ADDRESS"`

	NFTAddress string `mapstructure:"NFT_ADDRESS"`

	WalletNFT string `mapstructure:"WALLET_NFT"`

	WalletDTrash   string `mapstructure:"WALLET_DTRASH"`
	WalletLiquidez string `mapstructure:"WALLET_LIQUIDEZ"`
	WalletUsuarios string `mapstructure:"WALLET_USUARIOS"`

	WalletDTrashPerc   float64 `mapstructure:"WALLET_DTRASH_PERC"`
	WalletLiquidezPerc float64 `mapstructure:"WALLET_LIQUIDEZ_PERC"`
	WalletUsuariosPerc float64 `mapstructure:"WALLET_USUARIOS_PERC"`
	WalletRecyclePerc  float64 `mapstructure:"WALLET_RECYCLE_PERC"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
