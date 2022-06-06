import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlDeleteApiButton from '../../../../Components/Features/Graphql/api-instance/api-details/GraphqlDeleteApiButton';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title: `Graphql / api-instance / api-details / ${GraphqlDeleteApiButton.name}`,
  component: GraphqlDeleteApiButton,
} as unknown as ComponentMeta<typeof GraphqlDeleteApiButton>;

const apiRef = new ClusterObjectRef();

const Template: ComponentStory<typeof GraphqlDeleteApiButton> = args => (
  <MemoryRouter>
    <GraphqlDeleteApiButton {...args} />
  </MemoryRouter>
);

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = { apiRef: apiRef.toObject() } as Partial<
  typeof GraphqlDeleteApiButton
>;
