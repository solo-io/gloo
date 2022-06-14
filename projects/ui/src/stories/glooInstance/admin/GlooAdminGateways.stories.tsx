import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminGateways } from '../../../Components/Features/GlooInstance/Admin/GlooAdminGateways';

// TODO:  Add in mock from jest
export default {
  title: 'GlooInstance / Admin / GlooAdminGateways',
  component: GlooAdminGateways,
} as ComponentMeta<typeof GlooAdminGateways>;

const Template: ComponentStory<typeof GlooAdminGateways> = args => (
  // @ts-ignore
  <GlooAdminGateways {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminGateways>;
