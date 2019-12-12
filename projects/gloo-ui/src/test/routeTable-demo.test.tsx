import puppeteer, { Page, Browser } from 'puppeteer';

const testValues = {
  domain: 'http://localhost:3000/',
  routeName: '/vets',
  demoUrl: 'http://valet-petclinic-f3f885f59f.corp.solo.io'
};

const testVirtualService1 = {
  name: 'test-virtual-service-glooui',
  routeName: '/',
  upstream: 'gloo-system-apiserver-ui-8080'
};

const testVirtualService2 = {
  name: 'test-virtual-service-extauth',
  route1: {
    routeName: '/',
    upstream: 'extauth'
  },
  route2: {
    routeName: '/testing'
  }
};

const testRouteTable = {
  name: 'test-route-table',
  routeName: '/testing/test',
  upstreamName: 'default-kubernetes-443'
};

describe('route table demo', () => {
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

  describe('demo', () => {
    test('Overview Page renders correctly', async () => {
      await expect(page).toMatch('Enterprise Gloo Overview');

      await expect(page).toClick(
        'div[data-testid="view-virtual-services-link"]'
      );
    });

    test('create test virtual services 1', async () => {
      await expect(page).toClick('a[data-testid="virtual-services-navlink"]');

      await expect(page).toClick(
        'div[data-testid="create-virtual-service-modal"]'
      );
    }, 9999);

    test('Enter virtual service name', async () => {
      await expect(page).toFill(
        'input[name="virtualServiceName"]',
        testVirtualService1.name
      );
    }, 9999);

    test('submit create virtual service form ', async () => {
      await expect(page).toClick('button', { text: 'Create Virtual Service' });
    }, 9999);

    test('Add test routes ', async () => {
      await page.goto(
        `http://localhost:3000/virtualservices/gloo-system/${testVirtualService1.name}`
      );
      await expect(page).toMatch(testVirtualService1.name);
      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    });

    test('Enter route path', async () => {
      await expect(page).toFill(
        'input[name="path"]',
        testVirtualService1.routeName
      );
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testVirtualService1.upstream,
        timeout: 0
      });
    });

    test('submit form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route' });
    }, 9999);

    test('create test virtual services 2', async () => {
      await expect(page).toClick('a[data-testid="virtual-services-navlink"]');

      await expect(page).toClick(
        'div[data-testid="create-virtual-service-modal"]'
      );
    }, 9999);

    test('Enter virtual service name', async () => {
      await expect(page).toFill(
        'input[name="virtualServiceName"]',
        testVirtualService2.name
      );
    }, 9999);

    test('submit create virtual service form ', async () => {
      await expect(page).toClick('button', { text: 'Create Virtual Service' });
    }, 9999);

    test(' Add test routes ', async () => {
      await page.goto(
        `http://localhost:3000/virtualservices/gloo-system/${testVirtualService2.name}`
      );
      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    });

    test('Enter route path', async () => {
      await expect(page).toFill(
        'input[name="path"]',
        testVirtualService2.route1.routeName
      );
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testVirtualService2.route1.upstream,
        timeout: 0
      });
    });

    test('submit form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route' });
    }, 9999);

    // creat
    test('Create Route Table form', async () => {
      await expect(page).toClick('a[data-testid="virtual-services-navlink"]');

      await expect(page).toClick('div[data-testid="create-route-table-modal"]');
    });

    test('Enter route table name', async () => {
      await expect(page).toFill(
        'input[name="routeTableName"]',
        testRouteTable.name
      );
    }, 9999);

    test('submit create route table form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route Table' });
    }, 9999);

    test('make sure route table was created', async () => {
      await expect(page).toClick('a[data-testid="virtual-services-navlink"]');
      await expect(page).toMatch(testRouteTable.name);
      await page.goto(
        `http://localhost:3000/routetables/gloo-system/${testRouteTable.name}`
      );
    }, 9999);

    test('create a route', async () => {
      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    });
    test('Enter route path', async () => {
      await expect(page).toFill('input[name="path"]', testRouteTable.routeName);
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testRouteTable.upstreamName,
        timeout: 0
      });
    });

    test('submit form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route' });
    }, 9999);

    test('route to route table', async () => {
      await page.goto(
        `http://localhost:3000/virtualservices/gloo-system/${testVirtualService2.name}`
      );

      await expect(page).toClick('div[data-testid="create-new-route-modal"]');
    });

    test('Enter route path', async () => {
      await expect(page).toFill(
        'input[name="path"]',
        testVirtualService2.route2.routeName
      );
    }, 9999);

    test('Select Route table', async () => {
      await expect(page).toClick('div[data-testid="destination-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: 'Route Table',
        timeout: 0
      });
    }, 9999);

    test('Select Upstream', async () => {
      await expect(page).toClick('div[data-testid="upstream-dropdown"]');

      await expect(page).toClick('li[role="option"]', {
        text: testRouteTable.name,
        timeout: 0
      });
    });

    test('submit form ', async () => {
      await expect(page).toClick('button', { text: 'Create Route' });
    }, 9999);
  });

  test('can delete resources without errors', async () => {
    await page.goto(`http://localhost:3000/virtualservices/`);
    await expect(page).toClick('div[data-testid="delete-button"]');
    await expect(page).toMatch(
      'Are you sure you want to delete this virtual service?'
    );
    await expect(page).toClick('button', { text: 'Yes' });

    await expect(page).toClick('div[data-testid="delete-button"]');
    await expect(page).toMatch(
      'Are you sure you want to delete this virtual service?'
    );
    await expect(page).toClick('button', { text: 'Yes' });

    await expect(page).toClick('div[data-testid="delete-button"]');
    await expect(page).toMatch(
      'Are you sure you want to delete this Route Table?'
    );
    await expect(page).toClick('button', { text: 'Yes' });
  }, 15000);
});
