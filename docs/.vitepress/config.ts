import { defineConfig } from 'vitepress'
import { generateSidebar } from 'vitepress-sidebar';
import umlPlugin from 'markdown-it-plantuml';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "DIDimo",
  description: "Master the complexities of SSI identity solutions with DIDimo: Your go-to platform for testing, validating, and ensuring compliance in the ever-evolving digital identity ecosystem.",
  base: "/DIDimo/",

  lastUpdated: true,
  metaChunk: true,

  head: [
    [
      'script',
      {},
      `window.$crisp=[];window.CRISP_WEBSITE_ID="8dd97823-ddac-401e-991a-7498234e4f00";(function(){d=document;s=d.createElement("script");s.src="https://client.crisp.chat/l.js";s.async=1;d.getElementsByTagName("head")[0].appendChild(s);})();`
    ]
  ],
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
    logo: "",
    socialLinks: [
      { icon: 'github', link: 'https://github.com/forkbombeu/DIDimo' },
      { icon: "linkedin", link: "https://linkedin.com/company/forkbomb" },
    ],

    footer: {
      message:
        'Released under the <a href="https://github.com/forkbombeu/didimo/blob/main/LICENSE">AGPLv3 License</a>.',
      copyright:
        'Copyleft 🄯 2024-present <a href="https://forkbomb.solutions">The Forkbomb company</a>',
    },
  },
  markdown: {
    config(md) {
      md.use(umlPlugin)
    }
  }
})
