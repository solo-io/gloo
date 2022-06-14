import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedGateways } from '../../Components/Features/Admin/FederatedGateways';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedGateways',
  component: FederatedGateways,
} as ComponentMeta<typeof FederatedGateways>;

const Template: ComponentStory<typeof FederatedGateways> = args => (
  // @ts-ignore
  <FederatedGateways {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedGateways>;
