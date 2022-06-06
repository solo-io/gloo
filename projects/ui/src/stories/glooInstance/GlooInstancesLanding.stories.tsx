import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstancesLanding } from '../../Components/Features/GlooInstance/GlooInstancesLanding';

export default {
  title: `GlooInstance / ${GlooInstancesLanding.name}`,
  component: GlooInstancesLanding,
} as unknown as ComponentMeta<typeof GlooInstancesLanding>;

const Template: ComponentStory<typeof GlooInstancesLanding> = args => (
  // @ts-ignore
  <GlooInstancesLanding {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstancesLanding>;
