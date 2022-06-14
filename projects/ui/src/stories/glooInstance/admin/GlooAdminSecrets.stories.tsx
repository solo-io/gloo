import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminSecrets } from '../../../Components/Features/GlooInstance/Admin/GlooAdminSecrets';

export default {
  title: 'GlooInstance / Admin / GlooAdminSecrets',
  component: GlooAdminSecrets,
} as ComponentMeta<typeof GlooAdminSecrets>;

const Template: ComponentStory<typeof GlooAdminSecrets> = args => (
  // @ts-ignore
  <GlooAdminSecrets {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminSecrets>;
