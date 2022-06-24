import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import GraphqlApiDetails from '../../../../Components/Features/Graphql/api-instance/api-details/GraphqlApiDetails';
import { MemoryRouter } from 'react-router';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';
import {
  createClusterObjectRef,
  createGraphqlApi,
  createProxy,
} from 'stories/mocks/generators';
import { graphqlConfigApi } from 'API/graphql';

export default {
  title: 'Graphql / api-instance / api-details / GraphqlApiDetails',
  component: GraphqlApiDetails,
} as ComponentMeta<typeof GraphqlApiDetails>;

const Template: ComponentStory<typeof GraphqlApiDetails> = (args: any) => {
  const usePageApiRefDi = injectable(Apis.usePageApiRef, () => {
    return {
      name: args.name,
      namespace: args.namespace,
      clusterName: args.clusterName,
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
  const useGetGraphqlApiYamlDi = injectable(Apis.useGetGraphqlApiYaml, () => {
    return {
      data: args.apiYaml,
      error: args.error,
      isValidating: args.isValidating,
      mutate: args.mutate,
    };
  });

  const getGraphqlApiDi = injectable(graphqlConfigApi.getGraphqlApi, () => {
    return Promise.resolve(args.data);
  });

  return (
    <DiProvider
      use={[
        usePageApiRefDi,
        useGetGraphqlApiDetailsDi,
        useGetGraphqlApiYamlDi,
        getGraphqlApiDi,
      ]}>
      <MemoryRouter>
        <GraphqlApiDetails />
      </MemoryRouter>
    </DiProvider>
  );
};

const objData = createClusterObjectRef();

const gqlData = createGraphqlApi();

const apiYaml = faker.random.words(15);

export const Primary = Template.bind({});
// @ts-ignore
Primary.args = {
  ...objData,
  data: gqlData,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
  apiYaml,
} as Partial<typeof GraphqlApiDetails>;
