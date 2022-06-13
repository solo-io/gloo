import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlEditApiButton from '../../../../../Components/Features/Graphql/api-instance/api-details/executable-api/GraphqlEditApiButton';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title: `Graphql / api-instance / api-details / executable-api / ${GraphqlEditApiButton.name}`,
  component: GraphqlEditApiButton,
} as unknown as ComponentMeta<typeof GraphqlEditApiButton>;

const Template: ComponentStory<typeof GraphqlEditApiButton> = args => (
  <MemoryRouter>
    <GraphqlEditApiButton {...args} />
  </MemoryRouter>
);

const apiRef = new ClusterObjectRef();

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
} as Partial<typeof GraphqlEditApiButton>;
