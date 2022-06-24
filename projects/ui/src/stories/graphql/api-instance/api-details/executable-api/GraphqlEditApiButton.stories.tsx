import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlEditApiButton from '../../../../../Components/Features/Graphql/api-instance/api-details/executable-api/GraphqlEditApiButton';
import { MemoryRouter } from 'react-router';
import { graphqlConfigApi } from 'API/graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createGraphqlApi } from 'stories/mocks/generators';

export default {
  title:
    'Graphql / api-instance / api-details / executable-api / GraphqlEditApiButton',
  component: GraphqlEditApiButton,
} as ComponentMeta<typeof GraphqlEditApiButton>;

const Template: ComponentStory<typeof GraphqlEditApiButton> = (args: any) => {
  const getGraphqlApiDi = injectable(graphqlConfigApi.getGraphqlApi, () => {
    return Promise.resolve(args.data);
  });
  return (
    <DiProvider use={[getGraphqlApiDi]}>
      <MemoryRouter>
        <GraphqlEditApiButton {...args} />
      </MemoryRouter>
    </DiProvider>
  );
};

const apiRef = new ClusterObjectRef();
const data = createGraphqlApi();
export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  apiRef: apiRef.toObject(),
  data,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof GraphqlEditApiButton>;
