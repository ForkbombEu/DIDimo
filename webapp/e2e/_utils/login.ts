// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { type Browser, expect, type Page } from '@playwright/test';

export async function login(page: Page, email: string, password: string) {
	await page.goto('/login');

	const emailField = page.getByPlaceholder('name@foundation.org');
	const passwordField = page.getByPlaceholder('•••••');
	await expect(emailField).toBeVisible();
	await expect(passwordField).toBeVisible();

	await emailField.fill(email);
	await passwordField.fill(password);

	const submitButton = page.getByRole('button', { name: 'Log in' });
	await expect(submitButton).toBeVisible();
	await submitButton.click();

	await expect(page).toHaveURL(/keypairoom/);
	if (/keypairoom\/regenerate/.test(page.url())) {
		await page.getByRole('link', { name: 'Forgot the "seed"? Regenerate it' }).click();
	}

	const firstQuestion = page.locator('input[name="questions\\.whereParentsMet"]');
	await expect(firstQuestion).toBeVisible();

	await firstQuestion.fill('p');
	await page.locator('input[name="questions\\.nameFirstPet"]').fill('p');
	await page.locator('input[name="questions\\.whereHomeTown"]').fill('p');
	await page.locator('input[name="questions\\.nameFirstTeacher"]').fill('p');
	await page.locator('input[name="questions\\.nameMotherMaid"]').fill('p');

	await page.getByRole('button', { name: 'Generate private keys' }).click();

	const dashboardButton = page.getByRole('link', { name: 'Go to Dashboard' });
	await expect(dashboardButton).toBeVisible();
	await dashboardButton.click();

	await expect(page).toHaveURL(/my/);
}

export async function userLogin(browser: Browser, index: 'A' | 'B' | 'C') {
	const email = process.env[`TEST_USER_${index}_MAIL`];
	const password = process.env[`TEST_USER_${index}_PASS`];

	if (!email || !password) throw new Error(`No test user ${index} email or password set in ENV`);

	const context = await browser.newContext();
	const page = await context.newPage();

	await login(page, email, password);

	return page;
}

export function randomId() {
	return (Math.random() + 1).toString(36).substring(7);
}
