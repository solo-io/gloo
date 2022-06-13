import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooInstancesDetails } from '../../Components/Features/GlooInstance/GlooInstanceDetails';
import { MemoryRouter } from 'react-router';

export default {
  title: `GlooInstance / ${GlooInstancesDetails.name}`,
  component: GlooInstancesDetails,
} as unknown as ComponentMeta<typeof GlooInstancesDetails>;

const Template: ComponentStory<typeof GlooInstancesDetails> = args => (
  <MemoryRouter><GlooInstancesDetails /></MemoryRouter>
);

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof GlooInstancesDetails>;
