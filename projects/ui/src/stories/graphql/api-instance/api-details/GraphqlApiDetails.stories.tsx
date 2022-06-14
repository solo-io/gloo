import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlApiDetails from '../../../../Components/Features/Graphql/api-instance/api-details/GraphqlApiDetails';
import { MemoryRouter } from 'react-router';

export default {
  title: 'Graphql / api-instance / api-details / GraphqlApiDetails',
  component: GraphqlApiDetails,
} as ComponentMeta<typeof GraphqlApiDetails>;

const Template: ComponentStory<typeof GraphqlApiDetails> = args => (
  <MemoryRouter>
    {/* @ts-ignore */}
    <GraphqlApiDetails {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {} as Partial<typeof GraphqlApiDetails>;
