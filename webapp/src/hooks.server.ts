// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// file initialized by the Paraglide-SvelteKit CLI - Feel free to edit it
import { sequence } from "@sveltejs/kit/hooks"
import { i18n } from "$lib/i18n"

// add your own hooks as part of the sequence here
export const handle = sequence(i18n.handle())
