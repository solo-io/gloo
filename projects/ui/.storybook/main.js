module.exports = {
  "stories": [
    // "../src/**/*.stories.mdx",
    "../src/stories/**/*.stories.@(js|jsx|ts|tsx|mdx)"
  ],
  "addons": [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "@storybook/addon-interactions",
    "@storybook/preset-create-react-app",
  ],
  "framework": "@storybook/react",
  "core": {
    "builder": "@storybook/builder-webpack5"
  },
  features: {
    breakingChangesV7: true,
    interactionsDebugger: true,
    // storybook test seems to use this setting to `true`, so it's a useful config for testing
    storyStoreV7: false,
  },
}
