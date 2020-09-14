package app

import (
	"fmt"

	"github.com/1321822099/pdf_create/app/service/cmd"
	"github.com/1321822099/pdf_create/app/utils/config"
)

func init() {
	onStart(config.LoadConfigs)
	onStart(cmd.InitPool)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}
