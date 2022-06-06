import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminSettings } from '../../../Components/Features/GlooInstance/Admin/GlooAdminSettings';

export default {
  title: `GlooInstance / Admin / ${GlooAdminSettings.name}`,
  component: GlooAdminSettings,
} as unknown as ComponentMeta<typeof GlooAdminSettings>;

const Template: ComponentStory<typeof GlooAdminSettings> = args => (
  // @ts-ignore
  <GlooAdminSettings {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminSettings>;
