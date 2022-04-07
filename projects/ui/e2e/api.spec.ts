/// <reference types="@types/jest" />
import puppeteer from 'puppeteer';
import isEqual from 'lodash/isEqual';

describe('API Page', () => {
    jest.setTimeout(15000);
    const goodName = 'goodname';
    let browser: puppeteer.Browser;
    let page: puppeteer.Page;
    beforeAll(async () => {
        browser = await puppeteer.launch({
            headless: process.env.HEADLESS !== 'false',
            devtools: process.env.HEADLESS === 'false',
            dumpio: process.env.VERBOSE === 'true',
            slowMo: Number(process.env.SLOWMO || 10),
        });
        page = await browser.newPage();
        if (process.env.VERBOSE === 'true') {
            page
                .on('console', message =>
                    console.log(`${message.type().substring(0, 3).toUpperCase()} ${message.text()}`))
                .on('pageerror', ({ message }) => console.log(message))
                .on('response', response =>
                    console.log(`${response.status()} ${response.url()}`))
                .on('requestfailed', request =>
                    console.log(`${request.failure().errorText} ${request.url()}`))
        }
    });
    afterAll(async () => {
        await browser.close();
    });
    it('Should be able to run jest', () => {
        expect(true).toBe(true);
    });
    // /apis page.
    describe('/apis', () => {
        it('Should be able to see the main apis page', async () => {
            const response = await page.goto('http://localhost:3000/apis/');
            expect(response.ok()).toBe(true);
        });

        describe('api create', () => {
            it('Should be able to see the create API button', async () => {
                await page.goto('http://localhost:3000/apis/', {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle2']
                });
                // [data-testid='graphql-table-action-download']
                const downloadButton = await page.waitForSelector(`[data-testid='landing-create-api']`);
                expect(downloadButton!.asElement()).not.toBeNull();
                await downloadButton!.click();
                const closeBtn = await page.waitForSelector(`[data-testid='solo-modal-close']`);
                await closeBtn!.click();
            });
        });

        describe('api download', () => {
            it('Should be able to see the download button', async () => {
                await page.goto('http://localhost:3000/apis/', {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle2']
                });
                const downloadButton = await page.waitForSelector(`[data-testid='graphql-table-action-download']`);
                expect(downloadButton!.asElement()).not.toBeNull();
            });
        });

        describe('api update', () => {
            it('Should be able to update the new api modal', async () => {
                await page.goto('http://localhost:3000/apis/', {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });
                const downloadButton = await page.waitForSelector(`[data-testid='landing-create-api']`);

                expect(downloadButton!.asElement()).not.toBeNull();
                await downloadButton!.click();

                const nameInput = await page.waitForSelector(`[data-testid='new-api-modal-name']`);

                await nameInput!.focus();
                await page.keyboard.type(goodName);
                const result = await page.$eval<string>(`[data-testid='new-api-modal-name']`, (el) => {
                    const newEl = el as HTMLInputElement;
                    return newEl.value;
                });
                expect(result).toBe(goodName);

                const radioInput = await page.waitForSelector(`[data-testid='new-api-modal-apitype']`);
                await radioInput!.focus();

                const fileInput = await page.waitForSelector(`[data-testid='new-api-modal-file-upload']`);

                await fileInput!.focus();
                const [fileChooser] = await Promise.all([
                    page.waitForFileChooser(),
                    page.click(`[data-testid='new-api-modal-file-upload-btn']`),
                ]);

                await fileChooser.accept([`${__dirname}/mocks/book.gql`]);

                const submitBtn = await page.waitForSelector(`[data-testid="new-api-modal-submit"]`);
                await submitBtn!.click();

            }, 20000);
        });

        describe(`/apis/gloo-instances/gloo-system/gloo/apis/gloo-system/${goodName}`, () => {
            it('Should be able to go to the goodName page', async () => {
                await page.goto(`http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${goodName}`, {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });
            }, 10000);
            it('Should be able to click on the resolver', async () => {
                await page.goto(`http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${goodName}`, {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });
                const resolver = await page.waitForSelector(`[data-testid="resolver-productsForHome"]`);
                await resolver!.click();
                const upstreamBtn = await page.waitForSelector(`[data-testid="upstream-tab"]`);
                await upstreamBtn!.click();

                const inputOption = await page.waitForSelector('.ant-select-selection-search-input');
                await inputOption!.click();

                // Select the first option in the list.
                const selectable = await page.waitForSelector('.ant-select-item.ant-select-item-option');
                await page.hover('.ant-select-item.ant-select-item-option');
                await selectable!.click();

                // The onChange event does not register until you click off.
                const upstreamSection = await page.waitForSelector(`[data-testid="upstream-section"]`);
                await upstreamSection!.click();

                const resolverTabConfig = await page.waitForSelector(`[data-testid="resolver-config-tab"]`);
                await resolverTabConfig!.click();

                const submitBtn = await page.waitForSelector(`[data-testid="resolver-wizard-submit"]`);
                await submitBtn!.click();

                // Verify a new resolverConfig has been created.
                const routeIcon = await page.waitForSelector(`[data-testid="route-productsForHome"]`);
                await routeIcon!.click();

            }, 20000);

            it('Should be able to verify the resolver values.', async () => {
                await page.goto(`http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${goodName}`, {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });
                // Verify a new resolverConfig has been created.
                const routeIcon = await page.waitForSelector(`[data-testid="route-productsForHome"]`);
                await routeIcon!.click();

                const resolverTabConfig = await page.waitForSelector(`[data-testid="resolver-config-tab"]`);
                await resolverTabConfig!.click();

                await page.waitForSelector('#resolverConfiguration > textarea');

                const textAreaValues = await page.$$eval('.ace_line', (elements) => {
                    return elements.map((el) => {
                        return el!.textContent!.trim();
                    }).filter(el => el);
                });

                const expectedArray = [
                    'request:',
                    'headers:',
                    ':method:',
                    ':path:',
                    'queryParams:',
                    'body:',
                    'response:',
                    'resultRoot:',
                    'setters:',
                ];

                expect(isEqual(textAreaValues, expectedArray)).toBe(true);

            }, 20000);
        });

        describe('delete resolverConfig', () => {
            it('Should be able to remove a resolver', async () => {
                await page.goto(`http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${goodName}`, {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });

                const resolver = await page.waitForSelector(`[data-testid="resolver-productsForHome"]`);
                await resolver!.click();

                const removeResolverBtn = await page.waitForSelector(`[data-testid="remove-configuration-btn"]`);
                await removeResolverBtn!.click();

                const confirmBtnDelete = await page.waitForSelector(`[data-testid="confirm-delete-resolver"]`);
                await confirmBtnDelete!.click();
            });
        });

        describe('api delete', () => {
            it('Should delete a file if it exists', async () => {
                await page.goto('http://localhost:3000/apis/', {
                    waitUntil: ['load', 'domcontentloaded', 'networkidle0', 'networkidle2']
                });
                const hasGoodName = await page.waitForSelector(`.${goodName}`);
                const deleteBtn = await hasGoodName!.waitForSelector(`[data-testid="graphql-table-action-delete"]`);
                await deleteBtn!.click();
                const confirmDelete = await page.waitForSelector(`[data-testid="confirm-delete-graphql-table-row"]`);
                await confirmDelete!.click();
            });
        });

    });



});
