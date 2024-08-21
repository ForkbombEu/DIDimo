import { defineConfig } from 'vitepress'
import { generateSidebar } from 'vitepress-sidebar';
import umlPlugin from 'markdown-it-plantuml';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "DIDimo",
  description: "Master the complexities of SSI identity solutions with DIDimo: Your go-to platform for testing, validating, and ensuring compliance in the ever-evolving digital identity ecosystem.",
  base: "/didimo/",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Examples', link: '/markdown-examples' }
    ],

    sidebar: generateSidebar({
      useTitleFromFileHeading: true,
      sortMenusOrderNumericallyFromLink: true,

    }),

    socialLinks: [
      { icon: 'github', link: 'https://github.com/vuejs/vitepress' }
    ]
  },
  markdown: {
    config(md) {
      md.use(umlPlugin)
    }
  }
})
