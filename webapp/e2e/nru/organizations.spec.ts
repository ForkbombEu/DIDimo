import { test, expect, type Page } from '@playwright/test';
import { userLogin, randomId } from '@utils/login';
import { config } from 'dotenv';
import { nanoid } from 'nanoid';

config();

test.describe('it should test organizations and members', () => {
	let page: Page;
	let orgName: string;
	let orgId: string;

	test.beforeAll(async ({ browser }) => {
		const context = await browser.newContext();
		page = await context.newPage();
		await page.goto('/my');
	});

	test.afterAll(async ({ browser }) => {
		await browser.close();
	});

	test('userA should create an organization', async ({ browser }) => {
		page = await userLogin(browser, 'A');

		await page.goto('/my');
		await expect(page).toHaveURL(/my/);

		const organizationsLink = page.getByRole('link', { name: 'My organizations' }).first();
		await expect(organizationsLink).toBeVisible();
		await organizationsLink.click();

		await page.getByRole('link', { name: 'Create a new organization' }).click();

		orgName = `org-${nanoid(5)}`;
		await page.locator('input[name="name"]').click();
		await page.locator('input[name="name"]').fill(orgName);
		await page.getByRole('button', { name: 'Create organization' }).click();

		const organizationTitle = page.getByRole('button', { name: orgName });
		await expect(organizationTitle).toBeVisible();

		orgId = page.url().split('/').at(-1);
	});

	test('it should edit organization name', async () => {
		await page.getByRole('tab', { name: 'Settings' }).click();
		await expect(page).toHaveURL(/my\/organizations\/[^/]+\/settings/);

		await page.locator('input[name="name"]').click();

		orgName = `org-${randomId()}`;
		await page.locator('input[name="name"]').fill(orgName);
		await page.getByRole('button', { name: 'Save changes' }).click();
		await expect(page.getByRole('heading', { name: orgName })).toBeVisible();
	});

	test('userB should not see the organization page', async ({ browser }) => {
		page.close();
		page = await userLogin(browser, 'B');
		await page.goto(`/my/organizations/${orgId}`);
		await expect(page.getByText('404').first()).toBeVisible();
		page.close();
	});

	test('user A should invite user B and C to the organization', async ({ browser }) => {
		page = await userLogin(browser, 'A');
		await page.goto(`/my/organizations/${orgId}`);

		await page.getByRole('tab', { name: 'Members' }).click();
		await expect(page).toHaveURL(/my\/organizations\/[^/]+\/members/);

		const userB = 'userB@example.org';
		const userC = 'userC@example.org';

		await page.getByRole('button', { name: 'invite_members' }).click();
		await page.locator('textarea').fill(`${userB}\n${userC}`);
		await page.getByRole('button', { name: 'Review and confirm' }).click();

		const requestPromise = page.waitForRequest(
			(req) => req.url().includes('organizations/invite') && req.method() == 'POST'
		);
		await page.getByRole('button', { name: 'Send invites' }).click();
		const request = await requestPromise;
		await request.response();

		await expect(page.getByRole('dialog')).toBeHidden();
	});

	test('user B should accept user A invite', async ({ browser }) => {
		page.close();
		page = await userLogin(browser, 'B');
		await acceptInvite(page, orgName);
	});

	test('user C should accept user A invite', async ({ browser }) => {
		page.close();
		page = await userLogin(browser, 'C');
		await acceptInvite(page, orgName);
	});

	test('user A should see user B and C memberships', async ({ browser }) => {
		page.close();
		page = await userLogin(browser, 'A');
		await page.goto(`/my/organizations/${orgId}/members`);
		await expect(page.getByText('userB')).toBeVisible();
		await expect(page.getByText('userC')).toBeVisible();
	});

	test('user A should appoint user B to admin position', async () => {
		await page.goto(`/my/organizations/${orgId}/members`);

		const editRoleButton = page
			.getByText('userB Edit role')
			.getByRole('button', { name: 'Edit role' });
		await expect(editRoleButton).toBeVisible();
		await editRoleButton.click();

		await page.getByLabel('Role *').click();
		await page.getByRole('option', { name: 'admin' }).click();
		await page.getByRole('button', { name: 'Edit record' }).click();
		await expect(page.getByRole('dialog')).toBeHidden();

		await expect(page.getByText('Admin')).toBeVisible();
	});

	test("routing should hide the 'settings' page to admin", async ({ browser, page }) => {
		page.close();
		page = await userLogin(browser, 'B');
		await page.goto('/my/organizations');

		await expect(page.getByRole('main').getByText(orgName)).toBeVisible();

		const settingsButton = page.getByTestId(`${orgName} link`);
		await expect(settingsButton).toBeHidden();

		await page.getByRole('main').getByRole('link', { name: orgName }).click();
		await expect(page.getByRole('tab', { name: 'Settings' })).toBeHidden();

		await page.goto(`/my/organizations/${orgId}/settings`);
		await expect(page.getByText('404').first()).toBeVisible();
	});

	test("routing should hide the 'settings' section to user", async ({ browser, page }) => {
		page.close();
		page = await userLogin(browser, 'C');
		await page.goto('/my/organizations');

		await expect(page.getByRole('main').getByText(orgName)).toBeVisible();

		const settingsButton = page.getByTestId(`${orgName} link`);
		await expect(settingsButton).toBeHidden();

		await page.getByRole('main').getByRole('link', { name: orgName }).click();
		await expect(page.getByRole('tab', { name: 'Settings' })).toBeHidden();
	});
});

async function acceptInvite(page: Page, orgName: string) {
	await page.goto('/my/organizations');

	await expect(page.getByText(orgName)).toBeVisible();
	await page.getByRole('button', { name: 'Accept invite' }).click();

	const requestPromise = page.waitForRequest((req) =>
		req.url().includes('collections/org_invites/records')
	);
	const request = await requestPromise;
	await request.response();

	await expect(page.getByRole('link', { name: orgName })).toBeVisible();
}
