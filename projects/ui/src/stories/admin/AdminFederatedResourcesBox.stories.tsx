import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminFederatedResourcesBox } from '../../Components/Features/Admin/AdminBoxSummary';

// TODO:  Add in mock from jest.

export default {
  title: `Admin / ${AdminFederatedResourcesBox.name}`,
  component: AdminFederatedResourcesBox,
} as unknown as ComponentMeta<typeof AdminFederatedResourcesBox>;

const Template: ComponentStory<typeof AdminFederatedResourcesBox> = args => (
  // @ts-ignore
  <AdminFederatedResourcesBox {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminFederatedResourcesBox>;
