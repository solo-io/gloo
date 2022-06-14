import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedSettingsTable } from '../../Components/Features/Admin/FederatedSettingsTable';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / FederatedSettingsTable',
  component: FederatedSettingsTable,
} as ComponentMeta<typeof FederatedSettingsTable>;

const Template: ComponentStory<typeof FederatedSettingsTable> = args => (
  <FederatedSettingsTable />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedSettingsTable>;
