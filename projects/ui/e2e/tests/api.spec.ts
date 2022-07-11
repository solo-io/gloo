/// <reference types="@types/jest" />
import puppeteer from 'puppeteer';
import {
  addSubGraph,
  createApi,
  deleteApi,
  removeSubGraph,
} from './utils/api-helpers';
import {
  initBrowser,
  screenshot,
  sleep,
  waitForAndClick,
  waitForSelectorToBeRemoved,
} from './utils/helpers';
import fs from 'fs';
import { isEqual } from 'lodash';

describe('API Tests', () => {
  //
  // --- GLOBALS --- //
  //
  jest.setTimeout(20000);
  let browser: puppeteer.Browser;
  let page: puppeteer.Page;
  //
  // --- JEST LIFECYCLE --- //
  //
  beforeAll(async () => {
    try {
      fs.rmdirSync(`${__dirname}/screenshots`, { recursive: true });
    } catch {}
    fs.mkdirSync(`${__dirname}/screenshots`);
    const { browser: newBrowser, page: newPage } = await initBrowser();
    browser = newBrowser;
    page = newPage;
  });
  afterAll(async () => {
    await browser.close();
  });

  const exeApiName = 'jest-exe-users';
  const exeApiUrl = `http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${exeApiName}`;
  const stitchedApiName = 'jest-stitched';
  const stitchedApiUrl = `http://localhost:3000/gloo-instances/gloo-system/gloo/apis/gloo-system/${stitchedApiName}`;
  // ======================================================================= //
  //
  // --- Setup --- //
  //
  describe('\n================ Setup ================', () => {
    // it('Should be able to see the main apis page', async () => {
    //   const response = await page.goto('http://localhost:3000/apis/');
    //   expect(response.ok()).toBe(true);
    // });
    describe('Create Executable APIs', () => {
      it('Should be able to create executable APIs', async () => {
        await createApi(page, 'jest-exe-book', 'executable', 'book.gql');
        await createApi(page, 'jest-exe-users', 'executable', 'users.gql');
        await createApi(
          page,
          'jest-exe-products',
          'executable',
          'products.gql'
        );
      });
    });

    describe('Create Stitched API', () => {
      it('Should be able to create a stitched API', async () => {
        await createApi(page, 'jest-stitched', 'stitched');
      });
    });

    describe('Download API', () => {
      it('Should be able to see a download API button', async () => {
        await page.goto('http://localhost:3000/apis/', {
          waitUntil: ['networkidle0'],
        });
        const downloadButton = await page.waitForSelector(
          `[data-testid='graphql-table-action-download']`
        );
        expect(downloadButton!.asElement()).not.toBeNull();
      });
    });
  });

  // ======================================================================= //
  describe('\n================ Actions ================', () => {
    //
    // --- Stitching --- //
    //
    describe('--- Stitched APIs ---', () => {
      describe('Sub Graphs', () => {
        it('Should be able to add sub graphs', async () => {
          await page.goto(stitchedApiUrl, { waitUntil: ['networkidle0'] });
          await addSubGraph(page, 'jest-exe-book');
          await addSubGraph(page, 'jest-exe-products');
        });
        it('Should be able to add sub graphs with type merge configs', async () => {
          await page.goto(stitchedApiUrl, { waitUntil: ['networkidle0'] });
          await addSubGraph(page, 'jest-exe-users', {
            User: `argsMap:\n  id: UserSearch.username\nqueryName: GetUser\nselectionSet: '{username}'`,
          });
        });
        it('Should be able to remove sub graphs', async () => {
          await page.goto(stitchedApiUrl, { waitUntil: ['networkidle0'] });
          await removeSubGraph(page, 'jest-exe-book');
          await removeSubGraph(page, 'jest-exe-products');
        });
      });
    });

    // ======================================================================= //
    //
    // --- Executable --- //
    //
    describe('--- Executable APIs ---', () => {
      const fieldNameWithResolveDirective = 'GetUser';
      describe('Resolve Config', () => {
        it('Should be able to create a resolve definition', async () => {
          await page.goto(exeApiUrl, { waitUntil: ['networkidle0'] });
          const simpleScreenshot = () =>
            screenshot(page, `AddResolveDirective`);
          await simpleScreenshot();
          //
          // Open @resolve configuration wizard.
          await waitForAndClick(
            page,
            `[data-testid="resolver-${fieldNameWithResolveDirective}"]`
          );
          await simpleScreenshot();
          //
          // Open the upstream tab
          await waitForAndClick(page, `[data-testid="upstream-tab"]`);
          // Select the first option in the list.
          await waitForAndClick(page, '.ant-select');
          await page.keyboard.press('ArrowDown');
          await page.keyboard.press('Enter');
          await page.keyboard.press('Tab');
          await simpleScreenshot();
          //
          // Open the resolve config tab
          await waitForAndClick(page, `[data-testid="resolver-config-tab"]`);
          await simpleScreenshot();
          const submitBtn = await page.waitForSelector(
            `[data-testid="resolver-wizard-submit"]`
          );
          await submitBtn.click();
          await simpleScreenshot();
          // Verify a new resolverConfig has been created.
          const resolverCreatedEl = await page.waitForSelector(
            `[data-testid="route-${fieldNameWithResolveDirective}"]`
          );
          expect(!!resolverCreatedEl).toBeTruthy();
          await simpleScreenshot();
        });

        it('Should be able to verify the resolution values.', async () => {
          await page.goto(exeApiUrl, { waitUntil: ['networkidle0'] });
          const simpleScreenshot = () =>
            screenshot(page, `VerifyResolveDirective`);
          await simpleScreenshot();
          //
          // Verify a new resolverConfig has been created.
          // Open the @resolve config modal and the resolver config tab.
          await waitForAndClick(
            page,
            `[data-testid="resolver-${fieldNameWithResolveDirective}"]`
          );
          await waitForAndClick(page, `[data-testid="resolver-config-tab"]`);
          await page.waitForSelector('#resolverConfiguration > textarea');
          //
          // Verify the text content of the config
          const textAreaValues = await page.$$eval('.ace_line', elements => {
            return elements
              .map(el => {
                return el!.textContent!.trim();
              })
              .filter(el => el);
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
          await simpleScreenshot();
          expect(isEqual(textAreaValues, expectedArray)).toBeTruthy();
          await waitForAndClick(page, `[data-testid="solo-modal-close"]`);
        });

        it('Should be able to remove a resolve definition', async () => {
          const simpleScreenshot = () => screenshot(page, `RemoveResolver`);
          await page.goto(exeApiUrl, { waitUntil: ['networkidle0'] });
          await simpleScreenshot();
          // Open @resolve configuration wizard.
          await waitForAndClick(
            page,
            `[data-testid="resolver-${fieldNameWithResolveDirective}"]`
          );
          // Remove.
          await waitForAndClick(
            page,
            `[data-testid="remove-configuration-btn"]`
          );
          //
          // The confirm-delete modal is already on the page, but is hidden.
          // so we have to wait for it to appear.
          await sleep(150);
          await simpleScreenshot();
          await waitForAndClick(page, `[data-testid="confirm-modal-button"]`);
          await sleep(150);
          //
          // Verify that the resolverConfig has been removed.
          await waitForSelectorToBeRemoved(
            page,
            `[data-testid="route-${fieldNameWithResolveDirective}"]`
          );
          await simpleScreenshot();
        });
      });
    });
  });

  describe('\n================ Teardown ================', () => {
    describe('Delete APIs', () => {
      it('Should be able to delete APIs', async () => {
        await deleteApi(page, 'jest-exe-book');
        await deleteApi(page, 'jest-exe-users');
        await deleteApi(page, 'jest-exe-products');
        await deleteApi(page, 'jest-stitched');
      });
    });
  });
});
