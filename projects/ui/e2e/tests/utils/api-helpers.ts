/// <reference types="@types/jest" />
import puppeteer from 'puppeteer';
import { screenshot, sleep, waitForAndClick, waitForSelectorToBeRemoved } from './helpers';

/**
 * Creates an API
 * @param apiName
 * @param apiType
 * @param exeSchemaName
 */
export const createApi = async (
  page: puppeteer.Page,
  apiName: string,
  apiType: 'executable' | 'stitched',
  exeSchemaName?: string
) => {
  const simpleScreenshot = () => screenshot(page, `CreateApi:${apiName}`);
  await page.goto('http://localhost:3000/apis/', { waitUntil: 'networkidle0' });
  //
  // Open the modal.
  const createApiBtn = await page.waitForSelector(`[data-testid='landing-create-api']`);
  expect(!!createApiBtn).toBeTruthy();
  await simpleScreenshot();
  await createApiBtn.click();
  //
  // Input the API name.
  const nameInput = await page.waitForSelector(`[data-testid='new-api-modal-name']`);
  await nameInput.focus();
  await page.keyboard.type(apiName);
  const result = await page.$eval<string>(`[data-testid='new-api-modal-name']`, el => (el as HTMLInputElement).value);
  expect(result).toBe(apiName);
  //
  // Set the API type.
  const radioInput = await page.waitForSelector(`[data-testid='new-api-modal-apitype']`);
  if (apiType === 'executable') {
    //
    // -- Executable
    const exeInput = await radioInput.$('div:nth-of-type(1)');
    expect(await exeInput.evaluate(e => e.textContent.toLowerCase())).toBe('executable');
    await exeInput.click();
    // Input the file
    const fileInput = await page.waitForSelector(`[data-testid='new-api-modal-file-upload']`);
    await fileInput.focus();
    const [fileChooser] = await Promise.all([
      page.waitForFileChooser(),
      page.click(`[data-testid='new-api-modal-file-upload-btn']`),
    ]);
    expect(!!exeSchemaName).toBeTruthy();
    await fileChooser.accept([`${__dirname}/../mocks/${exeSchemaName}`]);
  } else {
    //
    // -- Stitched
    const stitchedInput = await radioInput!.$('div:nth-of-type(2)');
    expect(await stitchedInput.evaluate(e => e.textContent.toLowerCase())).toBe('stitched');
    await stitchedInput.click();
  }
  await simpleScreenshot();
  //
  // Submit and create the API.
  await waitForAndClick(page, `[data-testid="new-api-modal-submit"]`);
  await sleep(100);
  const urlRegex = new RegExp(`http:\/\/localhost:3000\/gloo-instances\/.*\/.*\/apis\/.*\/${apiName}\/?`);
  expect(urlRegex.test(page.url())).toBeTruthy();
  await simpleScreenshot();
};

/**
 * Deletes an API
 * @param apiName
 */
export const deleteApi = async (page: puppeteer.Page, apiName: string) => {
  const simpleScreenshot = () => screenshot(page, `DeleteApi:${apiName}`);
  await page.goto('http://localhost:3000/apis/', { waitUntil: 'networkidle0' });
  //
  // Find the API to delete.
  const apiRow = await page.waitForSelector(`.${apiName}`);
  const deleteBtn = await apiRow.waitForSelector(`[data-testid="graphql-table-action-delete"]`);
  await simpleScreenshot();
  await deleteBtn.click();
  //
  // The confirm-delete modal is already on the page, but is hidden.
  // so we have to wait for it to appear.
  await sleep(150);
  await waitForAndClick(page, `[data-testid="confirm-delete-graphql-table-row"]`);
  await sleep(150);
  //
  // Wait for the API to disappear from the list
  await waitForSelectorToBeRemoved(page, `.${apiName}`);
  await simpleScreenshot();
};

/**
 * Adds a sub graph to a stitched API. Must be called from the stitched API's
 * API details page.
 * @param page
 * @param subGraphApiName
 * @param typeMergeMap
 */
export const addSubGraph = async (
  page: puppeteer.Page,
  subGraphApiName: string,
  typeMergeMap?: { [key: string]: string }
) => {
  const simpleScreenshot = () => screenshot(page, `AddSubGraph:${subGraphApiName}`);
  //
  // Expect to be on an API details page.
  const urlRegex = new RegExp(`http:\/\/localhost:3000\/gloo-instances\/.*\/.*\/apis\/.*\/.*\/?`);
  expect(urlRegex.test(page.url())).toBeTruthy();
  //
  // Click the add sub graph button.
  const addModalBtn = await page.waitForSelector(`[data-testid='add-sub-graph-modal-button']`);
  expect(!!addModalBtn).toBeTruthy();
  await simpleScreenshot();
  await addModalBtn.click();
  //
  // Search for + select the sub graph.
  const apiDropdownSearch = await waitForAndClick(page, '.ant-select');
  await page.keyboard.type(subGraphApiName);
  await page.keyboard.press('Enter');
  await page.keyboard.press('Tab');
  let dropdownResult = await apiDropdownSearch.evaluate(el => el.textContent.trim());
  await simpleScreenshot();
  expect(dropdownResult).toBe(subGraphApiName);
  //
  // Add the type merge configs if they ware supplied.
  if (!!typeMergeMap) {
    const typeNames = Object.keys(typeMergeMap);
    const typeNameDropdown = await page.waitForSelector(`[data-testid='type-merge-name-dropdown']`);
    const addTypeMergeBtn = await page.waitForSelector(`[data-testid='add-type-merge-button']`);
    for (let i = 0; i < typeNames.length; i++) {
      const typeName = typeNames[i];
      const typeMergeConfig = typeMergeMap[typeName];
      // Add the field type.
      await typeNameDropdown.click();
      await page.keyboard.type(typeName);
      await page.keyboard.press('Enter');
      await page.keyboard.press('Tab');
      dropdownResult = await typeNameDropdown.evaluate(el => el.textContent.trim());
      expect(dropdownResult).toBe(typeName);
      await addTypeMergeBtn.click();
      await simpleScreenshot();
      // Configure it.
      const typeMergeConfigTextbox = await waitForAndClick(page, `[data-testid='type-merge-${typeName}'] textarea`);
      await page.keyboard.down('Meta');
      await page.keyboard.press('A');
      await page.keyboard.up('Meta');
      await page.keyboard.press('Backspace');
      // Write the lines out, making sure the editor doesn't mess it up with auto-indents.
      const configTextLines = typeMergeConfig.split('\n');
      for (let j = 1; j < configTextLines.length; j++) await page.keyboard.press('Enter');
      for (let j = 1; j < configTextLines.length; j++) await page.keyboard.press('ArrowUp');
      for (let j = 0; j < configTextLines.length; j++) {
        await typeMergeConfigTextbox.type(configTextLines[j]);
        await page.keyboard.press('ArrowRight');
      }
      await simpleScreenshot();
    }
  }
  //
  // Add the sub graph (this exits the modal).
  await waitForAndClick(page, `[data-testid='add-sub-graph-button']`);
  await sleep(50);
  await simpleScreenshot();
};

/**
 * Removes a sub graph from a stitched API. Must be called from the stitched API's
 * API details page.
 * @param subGraphApiName
 */
export const removeSubGraph = async (page: puppeteer.Page, subGraphApiName: string) => {
  const simpleScreenshot = () => screenshot(page, `RemoveSubGraph:${subGraphApiName}`);
  //
  // Expect to be on an API details page.
  const urlRegex = new RegExp(`http:\/\/localhost:3000\/gloo-instances\/.*\/.*\/apis\/.*\/.*\/?`);
  expect(urlRegex.test(page.url())).toBeTruthy();
  //
  // Find the sub graph to remove.
  const actionsContainer = await page.waitForSelector(`.${subGraphApiName}-actions`);
  const removeBtn = await actionsContainer.waitForSelector(`[data-testid="remove-sub-graph"]`);
  await simpleScreenshot();
  await removeBtn.click();
  //
  // The confirm-delete modal is already on the page, but is hidden.
  // so we have to wait for it to appear.
  await sleep(150);
  await waitForAndClick(page, `[data-testid="confirm-remove-sub-graph"]`);
  await sleep(150);
  //
  // Wait for the API to disappear from the list
  await waitForSelectorToBeRemoved(page, `.${subGraphApiName}-actions`);
  await simpleScreenshot();
};
