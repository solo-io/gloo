import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstanceAdministration } from '../../../Components/Features/GlooInstance/Admin/GlooInstanceAdministration';
import { MemoryRouter } from 'react-router';

export default {
  title: 'GlooInstance / Admin / GlooInstanceAdministration',
  component: GlooInstanceAdministration,
} as ComponentMeta<typeof GlooInstanceAdministration>;

const Template: ComponentStory<typeof GlooInstanceAdministration> = args => (
  <MemoryRouter>
    <GlooInstanceAdministration />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstanceAdministration>;
