import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminInnerPagesWrapper } from '../../../Components/Features/GlooInstance/Admin/GlooAdminInnerPagesWrapper';
import { MemoryRouter } from 'react-router';

// TODO:  Add in mock from jest
export default {
  title: 'GlooInstance / Admin / GlooAdminInnerPagesWrapper',
  component: GlooAdminInnerPagesWrapper,
} as ComponentMeta<typeof GlooAdminInnerPagesWrapper>;

const Template: ComponentStory<typeof GlooAdminInnerPagesWrapper> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GlooAdminInnerPagesWrapper {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooAdminInnerPagesWrapper>;
