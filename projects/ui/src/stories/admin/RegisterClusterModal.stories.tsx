import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RegisterClusterModal } from '../../Components/Features/Admin/RegisterClusterModal';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / RegisterClusterModal',
  component: RegisterClusterModal,
} as ComponentMeta<typeof RegisterClusterModal>;

const Template: ComponentStory<typeof RegisterClusterModal> = args => (
  // @ts-ignore
  <RegisterClusterModal {...args} />
);

export const Primary = Template.bind({});
Primary.args = {
  modalOpen: true,
} as Partial<typeof RegisterClusterModal>;
