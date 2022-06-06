import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminClustersBox } from '../../Components/Features/Admin/AdminBoxSummary';

// TODO:  Add in mock from jest
export default {
  title: `Admin / ${AdminClustersBox.name}`,
  component: AdminClustersBox,
} as unknown as ComponentMeta<typeof AdminClustersBox>;

const Template: ComponentStory<typeof AdminClustersBox> = args => (
  // @ts-ignore
  <AdminClustersBox {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminClustersBox>;
