import puppeteer, { Page, Browser } from 'puppeteer';

const testValues = {
  domain: 'http://localhost:3000/',
  routeName: '/vets',
  demoUrl: 'http://valet-petclinic-f3f885f59f.corp.solo.io'
};

const testUpstreamStatic = {
  name: 'my-test-static-upstream',
  type: 'Static',
  hostAddr: 'jsonplaceholder.typicode.com',
  hostPort: '80'
};

const testUpstreamAWS = {
  name: 'aws-upstream-test',
  type: 'AWS',
  region: 'us-east-1',
  secret: {
    name: 'aws-cert',
    namespace: 'gloo-system',
    accessKey: 'access-key',
    secretKey: 'secret-key'
  }
};

const testUpstream3 = {
  name: 'azure-upstream-test',
  type: 'Azure'
};

const testUpstream4 = {
  name: 'consul-upstream-test',
  type: 'Consul'
};

describe('basic upstream functionality', () => {
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

  describe('static upstream', () => {
    test('Overview Page renders correctly', async () => {
      await expect(page).toMatch('Enterprise Gloo Overview');

      await expect(page).toClick('div[data-testid="view-upstreams-link"]');
    });

    test('can open create upstream modal', async () => {
      await expect(page).toClick('a[data-testid="upstreams-navlink"]');

      await expect(page).toClick('div[data-testid="create-upstream-modal"]');
    }, 9999);

    test('Enter upstream name', async () => {
      await expect(page).toFill('input[name="name"]', testUpstreamStatic.name);
      await expect(page).toClick('div[data-testid="upstream-type"]');
      await expect(page).toClick('li[role="option"]', {
        text: 'Static',
        timeout: 0
      });
    }, 9999);

    test('type specific fields are shown', async () => {
      await expect(page).toMatch('Static Upstream Settings');
      await expect(page).toFill(
        'input[placeholder="address..."]',
        testUpstreamStatic.hostAddr
      );
      await expect(page).toFill(
        'input[placeholder="port..."]',
        testUpstreamStatic.hostPort
      );

      await expect(page).toClick('div[data-testid="green-plus-button"]');
    }, 9999);

    test('submit create upstream form ', async () => {
      await expect(page).toClick('button', { text: 'Create Upstream' });
    }, 9999);
  });

  describe('aws upstreams', () => {
    test('navigate to secrets page ', async () => {
      await page.goto('http://localhost:3000/settings/secrets/');
      expect(page).toMatch('AWS Secrets');
    });

    test('can fill out aws secret values', async () => {
      await expect(page).toFill(
        'input[data-testid="1-secret-name"]',
        testUpstreamAWS.secret.name
      );
      await expect(page).toClick('div[data-testid="1-secret-namespace"]');

      await expect(page).toClick('li[role="option"]', {
        text: 'gloo-system',
        timeout: 0
      });

      await expect(page).toFill(
        'input[data-testid="aws-secret-accessKey"]',
        testUpstreamAWS.secret.accessKey
      );

      await expect(page).toFill(
        'input[data-testid="aws-secret-secretKey"]',
        testUpstreamAWS.secret.secretKey
      );

      await expect(page).toClick('svg[data-testid="1-secret-green-plus"]');
    }, 9999);

    test('aws secret was created correctly', async () => {
      await page.goto('http://localhost:3000/settings/secrets/');

      await expect(page).toMatch(testUpstreamAWS.secret.name);
    }, 9999);

    describe('create AWS upstream', () => {
      test('fill out aws upstream values', async () => {
        await page.goto('http://localhost:3000/upstreams/');
        await expect(page).toClick('div[data-testid="create-upstream-modal"]');
        await expect(page).toFill('input[name="name"]', testUpstreamAWS.name);
        await expect(page).toClick('div[data-testid="upstream-type"]');
        await expect(page).toClick('li[role="option"]', {
          text: 'Aws',
          timeout: 0
        });

        await expect(page).toClick('div[data-testid="aws-region"]');
        await expect(page).toClick('li[role="option"]', {
          text: testUpstreamAWS.region,
          timeout: 0
        });

        await expect(page).toClick('div[data-testid="aws-secret"]');
        await expect(page).toClick('li[role="option"]', {
          text: testUpstreamAWS.secret.name,
          timeout: 0
        });

        await expect(page).toClick('button', { text: 'Create Upstream' });
      }, 9999);

      test('successfully created aws upstream', async () => {
        await page.goto('http://localhost:3000/upstreams/');
        await expect(page).toMatch(testUpstreamAWS.name);
      });
    });
  });

  describe('azure upstreams', () => {
    test('can create upstream', async () => {});
  });

  describe('consul upstreams', () => {
    test('can create upstream', async () => {});
  });

  describe('kube upstreams', () => {
    test('can create upstream', async () => {});
  });
});
