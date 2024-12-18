import { defineConfig } from '@playwright/test';

export const storageState = 'test-results/.auth/user.json';

export default defineConfig({
	webServer: {
		command: 'pnpm build && pnpm preview',
		port: 4173
	},

	testDir: 'e2e',

	projects: [
		{ name: 'setup', testMatch: /.*\.setup\.ts/ },
		{
			name: 'nruTests',
			testMatch: /nru\/.*\.spec\.ts/,
			use: {
				video: 'retain-on-failure'
			}
		},
		{
			name: 'loggedTests',
			testMatch: /logged\/.*\.spec\.ts/,
			use: {
				video: 'retain-on-failure',
				storageState
			},
			dependencies: ['setup']
		}
	]
});
