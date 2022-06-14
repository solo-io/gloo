import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceRouteTables } from '../../Components/Features/GlooInstance/GlooInstanceRouteTables';

export default {
  title: 'GlooInstance / GlooInstanceRouteTables',
  component: GlooInstanceRouteTables,
} as ComponentMeta<typeof GlooInstanceRouteTables>;

const Template: ComponentStory<typeof GlooInstanceRouteTables> = args => (
  // @ts-ignore
  <GlooInstanceRouteTables {...args} />
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceRouteTables>;
