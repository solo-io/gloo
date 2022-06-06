import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminProxy } from '../../../Components/Features/GlooInstance/Admin/GlooAdminProxy';

// TODO:  Add in mock from jest
export default {
  title: `GlooInstance / Admin / ${GlooAdminProxy.name}`,
  component: GlooAdminProxy,
} as unknown as ComponentMeta<typeof GlooAdminProxy>;

const Template: ComponentStory<typeof GlooAdminProxy> = args => (
  // @ts-ignore
  <GlooAdminProxy {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminProxy>;
