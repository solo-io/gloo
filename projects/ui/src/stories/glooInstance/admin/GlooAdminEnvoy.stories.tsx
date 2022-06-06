import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminEnvoy } from '../../../Components/Features/GlooInstance/Admin/GlooAdminEnvoy';

// TODO:  Add in mock from jest
export default {
  title: `GlooInstance / Admin / ${GlooAdminEnvoy.name}`,
  component: GlooAdminEnvoy,
} as unknown as ComponentMeta<typeof GlooAdminEnvoy>;

const Template: ComponentStory<typeof GlooAdminEnvoy> = args => (
  // @ts-ignore
  <GlooAdminEnvoy {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminEnvoy>;
