import { test, expect, type Page } from '@playwright/test';
import { userLogin, randomId } from '@utils/login';
import { config } from 'dotenv';

config();

test.describe.configure({ mode: 'serial' });

let page: Page;
let orgName: string;

test('it should create an organization', async ({ browser }) => {
	page = await userLogin(browser, 'A');
	await expect(page).toHaveURL('/my');

	await page.goto('/my/organizations');

	await page.getByRole('link', { name: 'plus Create a new organization' }).click();
	await expect(page).toHaveURL('/my/organizations/create');

	orgName = `org-${randomId()}`;
	await page.locator('input[name="name"]').click();
	await page.locator('input[name="name"]').fill(orgName);
	await page.getByRole('button', { name: 'Create organization' }).click();
	await expect(page).toHaveURL(new RegExp('/my/organizations/(.*)'));
	await expect(page.getByRole('heading', { name: orgName })).toBeVisible();
});

test('it should edit organization name', async () => {
	await page.getByRole('link', { name: 'General' }).click();
	await page.locator('input[name="name"]').click();

	orgName = `org-${randomId()}`;
	await page.locator('input[name="name"]').fill(orgName);
	await page.getByRole('button', { name: 'Save changes' }).click();
	await expect(page.getByRole('heading', { name: orgName })).toBeVisible();
});

test('it should add user B to the organization as admin', async () => {
	let username = 'userB';
	let email = `${username}@example.org`;
	let role = 'admin';
	await page.getByRole('link', { name: 'Members' }).click();
	await page.getByRole('button', { name: 'plus Add entry' }).click();
	await page.locator('form div').filter({ hasText: 'Select' }).nth(3).click();
	await page.getByPlaceholder('Select').fill(username);
	await page.getByText(email).click();
	await page.locator('select[name="role"]').selectOption({ label: role });
	await page.getByRole('button', { name: 'Create record' }).click();
	await expect(
		page
			.locator('tr')
			.filter({ has: page.getByText(email) })
			.filter({ has: page.getByText(role) })
	).toBeVisible();
});

test('it should add user C to the organization as member', async () => {
	let username = 'userC';
	let email = `${username}@example.org`;
	let role = 'member';
	await page.getByRole('button', { name: 'plus Add entry' }).click();
	await page.locator('form div').filter({ hasText: 'Select' }).nth(3).click();
	await page.getByPlaceholder('Select').fill(username);
	await page.getByText(email).click();
	await page.locator('select[name="role"]').selectOption({ label: role });
	await page.getByRole('button', { name: 'Create record' }).click();
	await expect(
		page
			.locator('tr')
			.filter({ has: page.getByText(email) })
			.filter({ has: page.getByText(role) })
	).toBeVisible();
});

test('it should fail to add A as member', async () => {
	let username = 'userA';
	let email = `${username}@example.org`;
	let role = 'member';
	await page.getByRole('button', { name: 'plus Add entry' }).click();
	await page.locator('form div').filter({ hasText: 'Select' }).nth(3).click();
	await page.getByPlaceholder('Select').fill(username);
	await page.getByRole('dialog').getByRole('document').getByText(email).click();
	await page.locator('select[name="role"]').selectOption({ label: role });
	await page.getByRole('button', { name: 'Create record' }).click();
	await expect(page.getByText('Failed to create record.')).toBeVisible();
});

test("it should hide the 'general' section to admin", async ({ browser }) => {
	page.close();
	page = await userLogin(browser, 'B');
	await page.goto('/my/organizations');

	await expect(page.getByText(orgName)).toBeVisible();

	const settingsButton = page.getByTestId(`${orgName} link`);
	await expect(settingsButton).toBeVisible();
	await settingsButton.click();

	await expect(page.getByText('general')).toBeHidden();
	await expect(page.getByText('members')).toBeVisible();

	const orgUrl = page.url();

	await page.goto(`${orgUrl}/general`);
	await expect(page.getByText('unauthorized')).toBeVisible();

	await page.goto(`${orgUrl}/members`);
	await expect(page.getByRole('button', { name: 'Add entry' })).toBeVisible();
});

test("it should hide the 'settings' section to user", async ({ browser }) => {
	page.close();
	page = await userLogin(browser, 'C');
	await page.goto('/my/organizations');

	await expect(page.getByText(orgName)).toBeVisible();

	const settingsButton = page.getByTestId(`${orgName} link`);
	await expect(settingsButton).toBeHidden();

	await expect(page.getByText('general')).toBeHidden();
	await expect(page.getByText('members')).toBeHidden();
});
