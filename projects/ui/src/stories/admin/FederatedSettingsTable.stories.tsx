import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FederatedSettingsTable } from '../../Components/Features/Admin/FederatedSettingsTable';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${FederatedSettingsTable.name}`,
  component: FederatedSettingsTable,
} as unknown as ComponentMeta<typeof FederatedSettingsTable>;

const Template: ComponentStory<typeof FederatedSettingsTable> = args => (
  // @ts-ignore
  <FederatedSettingsTable {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof FederatedSettingsTable>;
