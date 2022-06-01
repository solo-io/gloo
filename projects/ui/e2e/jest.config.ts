export default {
  preset: 'jest-puppeteer',
  testMatch: ['**/?(*.)+(spec|test).[t]s'],
  testPathIgnorePatterns: ['/node_modules/', 'dist'], //
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  transform: {
    '^.+\\.ts?$': 'ts-jest',
  },
  globalSetup: 'jest-environment-puppeteer/setup', // will be called once before all tests are executed
  globalTeardown: 'jest-environment-puppeteer/teardown',
  testEnvironment: 'jest-environment-puppeteer',
  moduleNameMapper: {
    'ace-builds': '<rootDir>/node_modules/ace-builds',
  },
};
