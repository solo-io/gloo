import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ExecutableGraphqlApiDetails } from '../../../../../Components/Features/Graphql/api-instance/api-details/executable-api/ExecutableGraphqlApiDetails';
import { MemoryRouter } from 'react-router';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { createGraphqlApi } from 'stories/mocks/generators';
import { graphqlConfigApi } from 'API/graphql';

export default {
  title:
    'Graphql / api-instance / api-details / executable-api / ExecutableGraphqlApiDetails',
  component: ExecutableGraphqlApiDetails,
} as ComponentMeta<typeof ExecutableGraphqlApiDetails>;

const Template: ComponentStory<typeof ExecutableGraphqlApiDetails> = (
  args: any
) => {
  const getGraphqlApiDi = injectable(graphqlConfigApi.getGraphqlApi, () => {
    return Promise.resolve(args.data);
  });
  const useGetConsoleOptionsDi = injectable(Apis.useGetConsoleOptions, () => {
    return args.consoleOptions;
  });
  // @ts-ignore
  const useListUpstreamsDi = injectable(Apis.useListUpstreams, () => {
    return {
      data: [],
      mutate: args.mutate,
      error: args.error,
      isValidating: args.isValidating,
    };
  });
  const useGetGraphqlApiDetailsDi = injectable(
    Apis.useGetGraphqlApiDetails,
    () => {
      return {
        data: args.data,
        error: args.error,
        isValidating: args.isValidating,
        mutate: args.mutate,
      };
    }
  );
  return (
    <DiProvider
      use={[
        useGetGraphqlApiDetailsDi,
        useListUpstreamsDi,
        useGetConsoleOptionsDi,
        getGraphqlApiDi,
      ]}>
      <MemoryRouter>
        <ExecutableGraphqlApiDetails {...args} />
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
  consoleOptions: {
    readonly: false,
    apiExplorerEnabled: true,
    errorMessage: '',
  },
} as Partial<typeof ExecutableGraphqlApiDetails>;
