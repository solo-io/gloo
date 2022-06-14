import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedVirtualServices } from '../../Components/Features/Admin/FederatedVirtualServices';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedVirtualServices',
  component: FederatedVirtualServices,
} as ComponentMeta<typeof FederatedVirtualServices>;

const Template: ComponentStory<typeof FederatedVirtualServices> = args => (
  // @ts-ignore
  <FederatedVirtualServices {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedVirtualServices>;
