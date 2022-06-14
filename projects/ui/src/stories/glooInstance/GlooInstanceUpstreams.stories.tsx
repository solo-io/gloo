import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceUpstreams } from '../../Components/Features/GlooInstance/GlooInstanceUpstreams';
import { MemoryRouter } from 'react-router';

export default {
  title: 'GlooInstance / GlooInstanceUpstreams',
  component: GlooInstanceUpstreams,
} as ComponentMeta<typeof GlooInstanceUpstreams>;

const Template: ComponentStory<typeof GlooInstanceUpstreams> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GlooInstanceUpstreams {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceUpstreams>;
