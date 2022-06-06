import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminInnerPagesWrapper } from '../../Components/Features/Admin/AdminInnerPagesWrapper';
import { MemoryRouter } from 'react-router';

export default {
  title: `Admin / ${AdminInnerPagesWrapper.name}`,
  component: AdminInnerPagesWrapper,
} as unknown as ComponentMeta<typeof AdminInnerPagesWrapper>;

const Template: ComponentStory<typeof AdminInnerPagesWrapper> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <AdminInnerPagesWrapper {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {
  cardName: AdminInnerPagesWrapper.name,
} as Partial<typeof AdminInnerPagesWrapper>;
