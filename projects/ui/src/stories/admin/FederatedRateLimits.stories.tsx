import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedRateLimits } from '../../Components/Features/Admin/FederatedRateLimits';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedRateLimits',
  component: FederatedRateLimits,
} as ComponentMeta<typeof FederatedRateLimits>;

const Template: ComponentStory<typeof FederatedRateLimits> = args => (
  // @ts-ignore
  <FederatedRateLimits {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedRateLimits>;
