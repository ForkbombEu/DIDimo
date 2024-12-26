// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// https://vitepress.dev/guide/custom-theme
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import matomo from "@datagouv/vitepress-plugin-matomo";

export default {
  extends: DefaultTheme,
  enhanceApp: ({ app, router, siteData }) => {
    matomo({
      router: router,
      siteID: 10, // Replace with your site id
      trackerUrl: "https://matomo.dyne.org/" // Replace with your matomo url
    })
  }
} satisfies Theme

