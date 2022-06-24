import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RegisterClusterModal } from '../../Components/Features/Admin/RegisterClusterModal';
import { jest } from '@storybook/jest';

// TODO:  Add in mock from jest
export default {
  title: 'Admin / RegisterClusterModal',
  component: RegisterClusterModal,
} as ComponentMeta<typeof RegisterClusterModal>;

const Template: ComponentStory<typeof RegisterClusterModal> = args => (
  <RegisterClusterModal {...args} />
);

export const Primary = Template.bind({});
Primary.args = {
  modalOpen: true,
  onClose: jest.fn(() => {
    console.log('Closed modal');
  }),
} as Partial<typeof RegisterClusterModal>;
