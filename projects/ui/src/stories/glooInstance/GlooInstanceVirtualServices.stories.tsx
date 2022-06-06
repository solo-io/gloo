import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceVirtualServices } from '../../Components/Features/GlooInstance/GlooInstanceVirtualServices';
import { MemoryRouter } from 'react-router';

export default {
  title: `GlooInstance / ${GlooInstanceVirtualServices.name}`,
  component: GlooInstanceVirtualServices,
} as unknown as ComponentMeta<typeof GlooInstanceVirtualServices>;

const Template: ComponentStory<typeof GlooInstanceVirtualServices> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GlooInstanceVirtualServices {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceVirtualServices>;
