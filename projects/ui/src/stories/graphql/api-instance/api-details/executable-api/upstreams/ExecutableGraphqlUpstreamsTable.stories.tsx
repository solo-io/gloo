import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ExecutableGraphqlUpstreamsTable from '../../../../../../Components/Features/Graphql/api-instance/api-details/executable-api/upstreams/ExecutableGraphqlUpstreamsTable';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';

export default {
  title: `Graphql / api-instance / api-details / executable-api / upstreams / ${ExecutableGraphqlUpstreamsTable.name}`,
  component: ExecutableGraphqlUpstreamsTable,
} as unknown as ComponentMeta<typeof ExecutableGraphqlUpstreamsTable>;

const Template: ComponentStory<
  typeof ExecutableGraphqlUpstreamsTable
> = args => (
  <MemoryRouter>
    <ExecutableGraphqlUpstreamsTable {...args} />
  </MemoryRouter>
);

const apiRef = new ClusterObjectRef();

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
} as Partial<typeof ExecutableGraphqlUpstreamsTable>;
