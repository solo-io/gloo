import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { EnableGraphqlFeature } from '../../Components/Features/Graphql/EnableGraphqlFeature';
import { MemoryRouter } from 'react-router';

export default {
  title: 'Graphql / EnableGraphqlFeature',
  component: EnableGraphqlFeature,
} as ComponentMeta<typeof EnableGraphqlFeature>;

const Template: ComponentStory<typeof EnableGraphqlFeature> = args => (
  // @ts-ignore
  <MemoryRouter>
    <EnableGraphqlFeature {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  children: <div>Here we go.</div>,
} as Partial<typeof EnableGraphqlFeature>;
