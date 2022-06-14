import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GraphqlInstance } from '../../../Components/Features/Graphql/api-instance/GraphqlInstance';
import { MemoryRouter } from 'react-router';

export default {
  title: 'Graphql / api-instance / GraphqlInstance',
  component: GraphqlInstance,
} as ComponentMeta<typeof GraphqlInstance>;

const Template: ComponentStory<typeof GraphqlInstance> = args => (
  // @ts-ignore
  <MemoryRouter>
    <GraphqlInstance {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {} as Partial<typeof GraphqlInstance>;
