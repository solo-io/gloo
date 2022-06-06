import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedRouteTables } from '../../Components/Features/Admin/FederatedRouteTables';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${FederatedRouteTables.name}`,
  component: FederatedRouteTables,
} as unknown as ComponentMeta<typeof FederatedRouteTables>;

const Template: ComponentStory<typeof FederatedRouteTables> = args => (
  // @ts-ignore
  <FederatedRouteTables {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedRouteTables>;
