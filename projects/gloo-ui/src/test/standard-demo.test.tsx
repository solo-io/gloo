import puppeteer, { Page, Browser } from 'puppeteer';

const testValues = {
  domain: 'http://localhost:3000/',
  routeName: '/vets',
  upstreamName: 'default-kubernetes-443',
  demoUrl: 'http://valet-petclinic-f3f885f59f.corp.solo.io',
  routeTableName: 'my-route-table',
  virtualServiceName: 'test-virtual-service-standard'
};

describe('app', () => {
  let page: any;
  let browser: Browser;

  beforeAll(async () => {
    browser = await puppeteer.launch({
      headless: false,
      defaultViewport: {
        width: 2280,
        height: 850
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

  describe('standard demo path', () => {
    test('Overview Page', async () => {
      await expect(page).toMatch('Enterprise Gloo Overview');

      await expect(page).toClick(
        'div[data-testid="view-virtual-services-link"]'
      );
    }, 9999);

    test('create test virtual service', async () => {
      await expect(page).toClick('a[data-testid="virtual-services-navlink"]');

      await expect(page).toClick(
        'div[data-testid="create-virtual-service-modal"]'
      );
    }, 9999);

    test('Enter virtual service name', async () => {
      await expect(page).toFill(
        'input[name="virtualServiceName"]',
        testValues.virtualServiceName
      );
    }, 9999);

    test('submit create virtual service form ', async () => {
      await expect(page).toClick('button', { text: 'Create Virtual Service' });
    }, 9999);

    test('New Route Modal', async () => {
      await page.goto(
        `http://localhost:3000/virtualservices/gloo-system/${testValues.virtualServiceName}`
      );
      await expect(page).toMatch('Create Route');
      await expect(page).toMatch('View YAML Configuration');
      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    }, 9999);

    test('Enter route path', async () => {
      await expect(page).toFill('input[name="path"]', testValues.routeName);
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testValues.upstreamName,
        timeout: 0
      });
    }, 9999);

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
