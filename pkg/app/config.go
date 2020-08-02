package app

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/thanhbinhdoan1993/practice-kuard/pkg/keygen"
)

type Config struct {
	Debug        bool
	DebugRootDir string `mapstructure:"debug-sitedata-dir"`
	ServerAddr   string `mapstructure:"address"`
	TLSAddr      string `mapstructure:"tls-address"`
	TLSDir       string `mapstructure:"tls-dir"`

	KeyGen keygen.Config

	Liveness  debugprobe.ProbeConfig
	Readiness debugprobe.ProbeConfig
}

func (k *App) BindConfig(v *viper.Viper, fs *pflag.FlagSet) {
	// k.kg.
}
