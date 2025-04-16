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
	app.RootCmd.Short = "\033[38;2;255;100;0m .o88b. d8888b. d88888b d8888b. d888888b .88b  d88. d888888b \033[0m\n" +
		"\033[38;2;255;71;43md8P  Y8 88  `8D 88'     88  `8D   `88'   88'YbdP`88   `88'   \033[0m\n" +
		"\033[38;2;255;43;86m8P      88oobY' 88ooooo 88   88    88    88  88  88    88    \033[0m\n" +
		"\033[38;2;255;14;129m8b      88`8b   88~~~~~ 88   88    88    88  88  88    88    \033[0m\n" +
		"\033[38;2;236;0;157mY8b  d8 88 `88. 88.     88  .8D   .88.   88  88  88   .88.   \033[0m\n" +
		"\033[38;2;197;0;171m `Y88P' 88   YD Y88888P Y8888D' Y888888P YP  YP  YP Y888888P \033[0m\n" +
		"\033[38;2;159;0;186m                                                             \033[0m\n" +
		"                                 \033[48;2;0;0;139m\033[38;2;255;255;255m              :(){ :|:& };: \033[0m\n" +
		"                                 \033[48;2;0;0;139m\033[38;2;255;255;255m with ❤ by Forkbomb hackers \033[0m\n"

	routes.Setup(app)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
