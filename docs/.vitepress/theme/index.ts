// SPDX-FileCopyrightText: 2024 The Forkbomb Company
//
// SPDX-License-Identifier: AGPL-3.0-or-later

// https://vitepress.dev/guide/custom-theme
import type { Theme } from "vitepress";
import { h } from "vue";
import DefaultTheme from "vitepress/theme";
import matomo from "@datagouv/vitepress-plugin-matomo";
import CustomFooter from "./components/Footer.vue";
import { useData } from "vitepress";

export default {
    extends: DefaultTheme,
    Layout: () => {
        return h(DefaultTheme.Layout, null, {
            // https://vitepress.dev/guide/extending-default-theme#layout-slots
            "doc-after": () => h(CustomFooter),
            "layout-bottom": () => {
                const { page } = useData();
                const layout = page.value.frontmatter.layout;
                if (layout === "home" || layout === "page") {
                    return h(CustomFooter);
                }
                return null;
            },
        });
    },
    enhanceApp: ({ app, router, siteData }) => {
        matomo({
            router: router,
            siteID: 10, // Replace with your site id
            trackerUrl: "https://matomo.dyne.org/", // Replace with your matomo url
        });
    },
} satisfies Theme;
