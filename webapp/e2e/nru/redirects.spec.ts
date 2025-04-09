// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { test, expect } from '@playwright/test';

test(`redirects to "/login" if not logged in and in "/my"`, async ({ page }) => {
	await page.goto('/my');
	await expect(page).toHaveURL(/login/);
});
