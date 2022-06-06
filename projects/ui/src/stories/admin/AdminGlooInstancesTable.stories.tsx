import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminGlooInstancesTable } from '../../Components/Features/Admin/AdminGlooInstancesTable';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${AdminGlooInstancesTable.name}`,
  component: AdminGlooInstancesTable,
} as unknown as ComponentMeta<typeof AdminGlooInstancesTable>;

const Template: ComponentStory<typeof AdminGlooInstancesTable> = args => (
  // @ts-ignore
  <AdminGlooInstancesTable {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminGlooInstancesTable>;
