package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/thanhbinhdoan1993/practice-kuard/pkg/app"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/version"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	app := app.NewApp()

	v := viper.GetViper()

	app.BindConfig(v, pflag.CommandLine)

	pflag.Parse()

	log.Printf("Starting kuard version: %v", version.VERSION)
	log.Printf(strings.Repeat("*", 70))
	log.Printf("* WARNING: This server may expose sensitive")
	log.Printf("* and secret information. Be careful.")
	log.Printf(strings.Repeat("*", 70))

	dumpConfig(v)

	app.LoadConfig(v)
	app.Run()
}

func dumpConfig(v *viper.Viper) {
	b, err := json.MarshalIndent(v.AllSettings(), "", "  ")
	if err != nil {
		log.Printf("Could not dump config: %v", err)
		return
	}
	log.Printf("Config: \n%v\n", string(b))
}
