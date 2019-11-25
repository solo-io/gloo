import puppeteer, { Page, Browser } from 'puppeteer';
import expect from 'expect-puppeteer';

const testValues = {
  domain: 'http://localhost:3000/',
  routeName: '/vets',
  upstreamName: 'default-petclinic-vets-8080',
  demoUrl: 'http://valet-petclinic-f3f885f59f.corp.solo.io'
};

describe('app', () => {
  let page: any;
  let browser: Browser;

  beforeAll(async () => {
    browser = await puppeteer.launch({
      headless: false,
      defaultViewport: {
        width: 2035,
        height: 1306
      },
      slowMo: 100
    });
    page = await browser.newPage();
    await page.goto(testValues.domain);
  });

  afterAll(async () => {
    await page.close();
    await browser.close();
  });

  describe('demo path', () => {
    test('Overview Page', async () => {
      await expect(page).toMatch('Enterprise Gloo Overview');

      await expect(page).toClick(
        'div[data-testid="view-virtual-services-link"]'
      );
    });

    test('New Route Modal', async () => {
      await expect(page).toClick('div[data-testid="view-details-link"]');
      await expect(page).toMatch('Create Route');
      await expect(page).toMatch('View Raw Configuration');
      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    });

    test('Enter route path', async () => {
      await expect(page).toFill('input[name="path"]', testValues.routeName);
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testValues.upstreamName,
        timeout: 0
      });
    });

    test('submit form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route' });
    }, 9999);

    test('route was created', async () => {
      await expect(page).toMatch(testValues.routeName, { timeout: 0 });
    }, 9999);

    xtest('information was updated in demo site', async () => {
      await page.goto(testValues.demoUrl);
      await page.goto(`${testValues.demoUrl}${testValues.routeName}`);
      await page.reload();
      await expect(page).toMatch('City');
    });

    xtest('delete route afterwards', async () => {
      await page.goBack();
      await page.goBack();

      await expect(page).toClick('div', { text: 'x' });
      await expect(page).toMatch('div[role="tooltip"]');
      await expect(page).toClick('button', { text: 'Yes' });
    }, 9999);
  });
});
