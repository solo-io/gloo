import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Clusters } from '../../Components/Features/Admin/Clusters';
import { MemoryRouter } from 'react-router';

export default {
  title: 'Admin / Clusters',
  component: Clusters,
} as ComponentMeta<typeof Clusters>;

const Template: ComponentStory<typeof Clusters> = args => (
  <MemoryRouter>
    <Clusters />
  </MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof Clusters>;
