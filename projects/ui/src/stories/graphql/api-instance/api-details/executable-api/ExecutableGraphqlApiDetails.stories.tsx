import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ExecutableGraphqlApiDetails } from '../../../../../Components/Features/Graphql/api-instance/api-details/executable-api/ExecutableGraphqlApiDetails';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title:
    'Graphql / api-instance / api-details / executable-api / ExecutableGraphqlApiDetails',
  component: ExecutableGraphqlApiDetails,
} as ComponentMeta<typeof ExecutableGraphqlApiDetails>;

const Template: ComponentStory<typeof ExecutableGraphqlApiDetails> = args => (
  <MemoryRouter>
    <ExecutableGraphqlApiDetails {...args} />
  </MemoryRouter>
);

const apiRef = new ClusterObjectRef();

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
} as Partial<typeof ExecutableGraphqlApiDetails>;
