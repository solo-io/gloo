import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminLanding } from '../../Components/Features/Admin/AdminLanding';
import { MemoryRouter } from 'react-router';

export default {
  title: `Admin / ${AdminLanding.name}`,
  component: AdminLanding,
} as unknown as ComponentMeta<typeof AdminLanding>;

const Template: ComponentStory<typeof AdminLanding> = args => (
  <MemoryRouter><AdminLanding /></MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminLanding>;
