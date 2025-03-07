// .vitepress/config.ts
import { defineConfig } from "file:///Users/giovanniabbatepaolo/Documents/GitHub/DIDimo/docs/node_modules/vitepress/dist/node/index.js";
import { generateSidebar } from "file:///Users/giovanniabbatepaolo/Documents/GitHub/DIDimo/docs/node_modules/vitepress-sidebar/dist/index.js";
import umlPlugin from "file:///Users/giovanniabbatepaolo/Documents/GitHub/DIDimo/docs/node_modules/markdown-it-plantuml/index.js";
var config_default = defineConfig({
  title: "DIDimo",
  description: "Master the complexities of SSI identity solutions with DIDimo: Your go-to platform for testing, validating, and ensuring compliance in the ever-evolving digital identity ecosystem.",
  base: "/DIDimo/",
  lastUpdated: true,
  metaChunk: true,
  head: [
    [
      "script",
      {},
      `window.$crisp=[];window.CRISP_WEBSITE_ID="8dd97823-ddac-401e-991a-7498234e4f00";(function(){d=document;s=d.createElement("script");s.src="https://client.crisp.chat/l.js";s.async=1;d.getElementsByTagName("head")[0].appendChild(s);})();`
    ]
  ],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Get Started", link: "/Architecture/1_start.html" },
      { text: "API Reference", target: "_self", link: "/API/index.html" }
    ],
    sidebar: generateSidebar({
      useTitleFromFileHeading: true,
      sortMenusOrderNumericallyFromLink: true
    }),
    logo: "",
    socialLinks: [
      { icon: "github", link: "https://github.com/forkbombeu/DIDimo" },
      { icon: "linkedin", link: "https://linkedin.com/company/forkbomb" }
    ],
    footer: {
      message: 'Released under the <a href="https://github.com/forkbombeu/didimo/blob/main/LICENSE">AGPLv3 License</a>.',
      copyright: 'Copyleft \u{1F12F} 2024-present <a href="https://forkbomb.solutions">The Forkbomb company</a>'
    }
  },
  markdown: {
    config(md) {
      md.use(umlPlugin);
    }
  }
});
export {
  config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLnZpdGVwcmVzcy9jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvVXNlcnMvZ2lvdmFubmlhYmJhdGVwYW9sby9Eb2N1bWVudHMvR2l0SHViL0RJRGltby9kb2NzLy52aXRlcHJlc3NcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIi9Vc2Vycy9naW92YW5uaWFiYmF0ZXBhb2xvL0RvY3VtZW50cy9HaXRIdWIvRElEaW1vL2RvY3MvLnZpdGVwcmVzcy9jb25maWcudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL1VzZXJzL2dpb3Zhbm5pYWJiYXRlcGFvbG8vRG9jdW1lbnRzL0dpdEh1Yi9ESURpbW8vZG9jcy8udml0ZXByZXNzL2NvbmZpZy50c1wiOy8vIFNQRFgtRmlsZUNvcHlyaWdodFRleHQ6IDIwMjQgVGhlIEZvcmtib21iIENvbXBhbnlcbi8vXG4vLyBTUERYLUxpY2Vuc2UtSWRlbnRpZmllcjogQUdQTC0zLjAtb3ItbGF0ZXJcblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSBcInZpdGVwcmVzc1wiO1xuaW1wb3J0IHsgZ2VuZXJhdGVTaWRlYmFyIH0gZnJvbSBcInZpdGVwcmVzcy1zaWRlYmFyXCI7XG5pbXBvcnQgdW1sUGx1Z2luIGZyb20gXCJtYXJrZG93bi1pdC1wbGFudHVtbFwiO1xuXG4vLyBodHRwczovL3ZpdGVwcmVzcy5kZXYvcmVmZXJlbmNlL3NpdGUtY29uZmlnXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICB0aXRsZTogXCJESURpbW9cIixcbiAgZGVzY3JpcHRpb246XG4gICAgXCJNYXN0ZXIgdGhlIGNvbXBsZXhpdGllcyBvZiBTU0kgaWRlbnRpdHkgc29sdXRpb25zIHdpdGggRElEaW1vOiBZb3VyIGdvLXRvIHBsYXRmb3JtIGZvciB0ZXN0aW5nLCB2YWxpZGF0aW5nLCBhbmQgZW5zdXJpbmcgY29tcGxpYW5jZSBpbiB0aGUgZXZlci1ldm9sdmluZyBkaWdpdGFsIGlkZW50aXR5IGVjb3N5c3RlbS5cIixcbiAgYmFzZTogXCIvRElEaW1vL1wiLFxuXG4gIGxhc3RVcGRhdGVkOiB0cnVlLFxuICBtZXRhQ2h1bms6IHRydWUsXG5cbiAgaGVhZDogW1xuICAgIFtcbiAgICAgIFwic2NyaXB0XCIsXG4gICAgICB7fSxcbiAgICAgIGB3aW5kb3cuJGNyaXNwPVtdO3dpbmRvdy5DUklTUF9XRUJTSVRFX0lEPVwiOGRkOTc4MjMtZGRhYy00MDFlLTk5MWEtNzQ5ODIzNGU0ZjAwXCI7KGZ1bmN0aW9uKCl7ZD1kb2N1bWVudDtzPWQuY3JlYXRlRWxlbWVudChcInNjcmlwdFwiKTtzLnNyYz1cImh0dHBzOi8vY2xpZW50LmNyaXNwLmNoYXQvbC5qc1wiO3MuYXN5bmM9MTtkLmdldEVsZW1lbnRzQnlUYWdOYW1lKFwiaGVhZFwiKVswXS5hcHBlbmRDaGlsZChzKTt9KSgpO2AsXG4gICAgXSxcbiAgXSxcbiAgdGhlbWVDb25maWc6IHtcbiAgICAvLyBodHRwczovL3ZpdGVwcmVzcy5kZXYvcmVmZXJlbmNlL2RlZmF1bHQtdGhlbWUtY29uZmlnXG4gICAgbmF2OiBbXG4gICAgICB7IHRleHQ6IFwiSG9tZVwiLCBsaW5rOiBcIi9cIiB9LFxuICAgICAgeyB0ZXh0OiBcIkdldCBTdGFydGVkXCIsIGxpbms6IFwiL0FyY2hpdGVjdHVyZS8xX3N0YXJ0Lmh0bWxcIiB9LFxuICAgICAgeyB0ZXh0OiBcIkFQSSBSZWZlcmVuY2VcIiwgdGFyZ2V0OiBcIl9zZWxmXCIsIGxpbms6IFwiL0FQSS9pbmRleC5odG1sXCIgfSxcbiAgICBdLFxuXG4gICAgc2lkZWJhcjogZ2VuZXJhdGVTaWRlYmFyKHtcbiAgICAgIHVzZVRpdGxlRnJvbUZpbGVIZWFkaW5nOiB0cnVlLFxuICAgICAgc29ydE1lbnVzT3JkZXJOdW1lcmljYWxseUZyb21MaW5rOiB0cnVlLFxuICAgIH0pLFxuICAgIGxvZ286IFwiXCIsXG4gICAgc29jaWFsTGlua3M6IFtcbiAgICAgIHsgaWNvbjogXCJnaXRodWJcIiwgbGluazogXCJodHRwczovL2dpdGh1Yi5jb20vZm9ya2JvbWJldS9ESURpbW9cIiB9LFxuICAgICAgeyBpY29uOiBcImxpbmtlZGluXCIsIGxpbms6IFwiaHR0cHM6Ly9saW5rZWRpbi5jb20vY29tcGFueS9mb3JrYm9tYlwiIH0sXG4gICAgXSxcblxuICAgIGZvb3Rlcjoge1xuICAgICAgbWVzc2FnZTpcbiAgICAgICAgJ1JlbGVhc2VkIHVuZGVyIHRoZSA8YSBocmVmPVwiaHR0cHM6Ly9naXRodWIuY29tL2Zvcmtib21iZXUvZGlkaW1vL2Jsb2IvbWFpbi9MSUNFTlNFXCI+QUdQTHYzIExpY2Vuc2U8L2E+LicsXG4gICAgICBjb3B5cmlnaHQ6XG4gICAgICAgICdDb3B5bGVmdCBcdUQ4M0NcdUREMkYgMjAyNC1wcmVzZW50IDxhIGhyZWY9XCJodHRwczovL2Zvcmtib21iLnNvbHV0aW9uc1wiPlRoZSBGb3JrYm9tYiBjb21wYW55PC9hPicsXG4gICAgfSxcbiAgfSxcbiAgbWFya2Rvd246IHtcbiAgICBjb25maWcobWQpIHtcbiAgICAgIG1kLnVzZSh1bWxQbHVnaW4pO1xuICAgIH0sXG4gIH0sXG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFJQSxTQUFTLG9CQUFvQjtBQUM3QixTQUFTLHVCQUF1QjtBQUNoQyxPQUFPLGVBQWU7QUFHdEIsSUFBTyxpQkFBUSxhQUFhO0FBQUEsRUFDMUIsT0FBTztBQUFBLEVBQ1AsYUFDRTtBQUFBLEVBQ0YsTUFBTTtBQUFBLEVBRU4sYUFBYTtBQUFBLEVBQ2IsV0FBVztBQUFBLEVBRVgsTUFBTTtBQUFBLElBQ0o7QUFBQSxNQUNFO0FBQUEsTUFDQSxDQUFDO0FBQUEsTUFDRDtBQUFBLElBQ0Y7QUFBQSxFQUNGO0FBQUEsRUFDQSxhQUFhO0FBQUE7QUFBQSxJQUVYLEtBQUs7QUFBQSxNQUNILEVBQUUsTUFBTSxRQUFRLE1BQU0sSUFBSTtBQUFBLE1BQzFCLEVBQUUsTUFBTSxlQUFlLE1BQU0sNkJBQTZCO0FBQUEsTUFDMUQsRUFBRSxNQUFNLGlCQUFpQixRQUFRLFNBQVMsTUFBTSxrQkFBa0I7QUFBQSxJQUNwRTtBQUFBLElBRUEsU0FBUyxnQkFBZ0I7QUFBQSxNQUN2Qix5QkFBeUI7QUFBQSxNQUN6QixtQ0FBbUM7QUFBQSxJQUNyQyxDQUFDO0FBQUEsSUFDRCxNQUFNO0FBQUEsSUFDTixhQUFhO0FBQUEsTUFDWCxFQUFFLE1BQU0sVUFBVSxNQUFNLHVDQUF1QztBQUFBLE1BQy9ELEVBQUUsTUFBTSxZQUFZLE1BQU0sd0NBQXdDO0FBQUEsSUFDcEU7QUFBQSxJQUVBLFFBQVE7QUFBQSxNQUNOLFNBQ0U7QUFBQSxNQUNGLFdBQ0U7QUFBQSxJQUNKO0FBQUEsRUFDRjtBQUFBLEVBQ0EsVUFBVTtBQUFBLElBQ1IsT0FBTyxJQUFJO0FBQ1QsU0FBRyxJQUFJLFNBQVM7QUFBQSxJQUNsQjtBQUFBLEVBQ0Y7QUFDRixDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
