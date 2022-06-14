import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstancesLanding } from '../../Components/Features/GlooInstance/GlooInstancesLanding';
import { MemoryRouter } from 'react-router';

export default {
  title: 'GlooInstance / GlooInstancesLanding',
  component: GlooInstancesLanding,
} as ComponentMeta<typeof GlooInstancesLanding>;

const Template: ComponentStory<typeof GlooInstancesLanding> = args => (
  <MemoryRouter>
    <GlooInstancesLanding />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstancesLanding>;
