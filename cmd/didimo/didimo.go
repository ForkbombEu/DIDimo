// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// DIDimo is your companion tool for be compliant with your SSI system.
package main

import (
	"log"

	_ "github.com/forkbombeu/didimo/migrations"
	"github.com/forkbombeu/didimo/pkg/routes"

	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()
	app.RootCmd.Short = "\033[38;2;255;100;0m      dP oo       dP oo                     \033[0m\n" +
		"\033[38;2;255;71;43m      88          88                        \033[0m\n" +
		"\033[38;2;255;43;86m.d888b88 dP .d888b88 dP 88d8b.d8b. .d8888b. \033[0m\n" +
		"\033[38;2;255;14;129m88'  `88 88 88'  `88 88 88'`88'`88 88'  `88 \033[0m\n" +
		"\033[38;2;236;0;157m88.  .88 88 88.  .88 88 88  88  88 88.  .88 \033[0m\n" +
		"\033[38;2;197;0;171m`88888P8 dP `88888P8 dP dP  dP  dP `88888P' \033[0m\n" +
		"\033[38;2;159;0;186m                                             \033[0m\n" +
		"                   \033[48;2;0;0;139m\033[38;2;255;255;255m           :(){ :|:& };: \033[0m\n" +
		"                   \033[48;2;0;0;139m\033[38;2;255;255;255m by The Forkbomb Company \033[0m\n"

	routes.Setup(app)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
