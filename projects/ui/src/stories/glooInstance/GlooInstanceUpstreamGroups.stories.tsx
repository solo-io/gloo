import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceUpstreamGroups } from '../../Components/Features/GlooInstance/GlooInstanceUpstreamGroups';
import { MemoryRouter } from 'react-router';

export default {
  title: 'GlooInstance / GlooInstanceUpstreamGroups',
  component: GlooInstanceUpstreamGroups,
} as ComponentMeta<typeof GlooInstanceUpstreamGroups>;

const Template: ComponentStory<typeof GlooInstanceUpstreamGroups> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GlooInstanceUpstreamGroups {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceUpstreamGroups>;
