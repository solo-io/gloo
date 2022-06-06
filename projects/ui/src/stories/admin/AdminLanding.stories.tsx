import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminLanding } from '../../Components/Features/Admin/AdminLanding';

export default {
  title: `Admin / ${AdminLanding.name}`,
  component: AdminLanding,
} as unknown as ComponentMeta<typeof AdminLanding>;

const Template: ComponentStory<typeof AdminLanding> = args => (
  // @ts-ignore
  <AdminLanding {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminLanding>;
