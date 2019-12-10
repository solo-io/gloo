import puppeteer, { Page, Browser } from 'puppeteer';

const testValues = {
  domain: 'http://localhost:3000/',
  routeName: '/vets',
  demoUrl: 'http://valet-petclinic-f3f885f59f.corp.solo.io'
};

const testUpstream1 = {
  name: 'my-test-static-upstream',
  type: 'Static',
  hostAddr: 'jsonplaceholder.typicode.com',
  hostPort: '80'
};

xdescribe('YAML editing demo', () => {
  let page: Page;
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
    // close
    await page.close();
    await browser.close();
  });

  xdescribe('can edit yaml virtual services', () => {
    test('Overview Page renders correctly', async () => {
      await expect(page).toMatch('Enterprise Gloo Overview');

      await expect(page).toClick('div[data-testid="view-upstreams-link"]');
    });
  });
});
