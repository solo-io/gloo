import puppeteer, { Puppeteer } from 'puppeteer';

export const initBrowser = async () => {
  const browser = await puppeteer.launch({
    headless: process.env.HEADLESS !== 'false',
    devtools: process.env.HEADLESS === 'false',
    dumpio: process.env.VERBOSE === 'true',
    slowMo: Number(process.env.SLOWMO || 10),
  });
  const page = await browser.newPage();
  await page.setViewport({
    width: 1440,
    height: 1080,
  });
  if (process.env.VERBOSE === 'true') {
    page
      .on('console', message =>
        console.log(
          `${message.type().substring(0, 3).toUpperCase()} ${message.text()}`
        )
      )
      .on('pageerror', ({ message }) => console.log(message))
      .on('response', response =>
        console.log(`${response.status()} ${response.url()}`)
      )
      .on('requestfailed', request =>
        console.log(`${request.failure().errorText} ${request.url()}`)
      );
  }
  return { browser, page };
};

let numScreenshots = 0;
const makeScreenshotPath = (fileName: string) =>
  `${__dirname}/../screenshots/${(numScreenshots++)
    .toString()
    .padStart(3, '0')}_${fileName}.jpeg`;

export const screenshot = async (page: puppeteer.Page, fileName: string) =>
  await page.screenshot({ path: makeScreenshotPath(fileName) });

export const sleep = async (ms: number) =>
  await new Promise((resolve, _) => setTimeout(resolve, ms));

/**
 *  This will retry finding the selector until it errors when it isn't found,
 *  or until the test times out (there is probably a better way of doing this, but this works).
 * @param page
 * @param selector
 * @param interval
 */
export const waitForSelectorToBeRemoved = async (
  page: puppeteer.Page,
  selector: string,
  interval = 50
) => {
  while (true) {
    try {
      await page.waitForSelector(selector, { timeout: interval });
      await sleep(interval);
    } catch {
      break;
    }
  }
};

export const waitForAndClick = async (
  page: puppeteer.Page,
  selector: string
) => {
  const el = await page.waitForSelector(selector);
  await el.click();
  return el;
};
