import { test as setup } from '@playwright/test';
import { config } from 'dotenv';
import { userLogin } from '@utils/login';
import { storageState } from '../playwright.config';

config();

setup('authenticate', async ({ browser }) => {
	const page = await userLogin(browser, 'A');
	await page.context().storageState({ path: storageState });
});
