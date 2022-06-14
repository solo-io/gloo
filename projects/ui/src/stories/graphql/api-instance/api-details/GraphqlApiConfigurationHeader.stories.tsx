import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlApiConfigurationHeader from '../../../../Components/Features/Graphql/api-instance/api-details/GraphqlApiConfigurationHeader';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title: 'Graphql / api-instance / api-details / GraphqlApiConfigurationHeader',
  component: GraphqlApiConfigurationHeader,
} as ComponentMeta<typeof GraphqlApiConfigurationHeader>;

const Template: ComponentStory<typeof GraphqlApiConfigurationHeader> = args => (
  // @ts-ignore
  <MemoryRouter>
    <GraphqlApiConfigurationHeader {...args} />
  </MemoryRouter>
);

const apiRef = new ClusterObjectRef();

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
} as Partial<typeof GraphqlApiConfigurationHeader>;
