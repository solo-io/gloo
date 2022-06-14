import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminClustersBox } from '../../Components/Features/Admin/AdminBoxSummary';
import { MemoryRouter } from 'react-router';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { expect, jest } from '@storybook/jest';
import * as Apis from 'API/hooks';
import { within } from '@storybook/testing-library';
import { createClusterDetailsObj } from 'stories/mocks/generators';
import { ServiceError } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb_service';
import { grpc } from '@improbable-eng/grpc-web';
import { ClusterDetails } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';

export default {
  // Dynamic name modules generates errors when you run this in a test.
  // title: 'Admin / AdminClustersBox',
  component: AdminClustersBox,
} as ComponentMeta<typeof AdminClustersBox>;

interface TemplateType {
  clusterDetails?: ClusterDetails.AsObject[];
  error?: ServiceError;
  mutate?: any;
  isValidating?: boolean;
}

const Template: ComponentStory<typeof AdminClustersBox> = (
  args: TemplateType | any
) => {
  const useListClusterDetailsDi = injectable(Apis.useListClusterDetails, () => {
    return {
      data: args.clusterDetails,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  return (
    <DiProvider use={[useListClusterDetailsDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};

export const OneCluster = Template.bind({});
const clusterDetailsOne = Array.from({ length: 1 }).map(() => {
  return createClusterDetailsObj();
});
OneCluster.args = {
  clusterDetails: clusterDetailsOne,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof AdminClustersBox>;
OneCluster.parameters = {};
OneCluster.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-box-summary');
  expect(summary).not.toBeNull();
};

const TemplateTwo: ComponentStory<typeof AdminClustersBox> = (
  args: TemplateType | any
) => {
  const useListClusterDetailsDi = injectable(Apis.useListClusterDetails, () => {
    return {
      data: args.clusterDetails,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  return (
    <DiProvider use={[useListClusterDetailsDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};
export const TwoCluster = TemplateTwo.bind({});

const clusterDetailsTwo = Array.from({ length: 2 }).map(() => {
  return createClusterDetailsObj();
});

TwoCluster.args = {
  clusterDetails: clusterDetailsTwo,
  error: undefined,
  mutate: jest.fn(),
  isValidating: false,
} as Partial<typeof AdminClustersBox>;
TwoCluster.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-box-summary');
  expect(summary).not.toBeNull();
};
TwoCluster.parameters = {
  msw: {
    handlers: [],
  },
};

const TemplateError: ComponentStory<typeof AdminClustersBox> = (
  args: TemplateType | any
) => {
  const useListClusterDetailsDi = injectable(Apis.useListClusterDetails, () => {
    return {
      data: args.detailsCluster,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  return (
    <DiProvider use={[useListClusterDetailsDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};

const meta = new grpc.Metadata();
meta.append('thing-happened', 'A bad thing happened')!;
const error: ServiceError = {
  message: 'Some Error Happened!',
  code: 1,
  metadata: meta,
};

export const ErrorCluster = TemplateError.bind({});
ErrorCluster.args = {
  clusterDetails: [],
  error,
  isValidating: false,
  mutate: jest.fn(),
} as Partial<typeof AdminClustersBox>;
ErrorCluster.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const dataError = await canvas.getByTestId('data-error');
  expect(dataError).not.toBeNull();
};

ErrorCluster.parameters = {
  msw: {
    handlers: [],
  },
};

const TemplateLoading: ComponentStory<typeof AdminClustersBox> = (
  args: TemplateType | any
) => {
  const useListClusterDetailsData = injectable(
    Apis.useListClusterDetails,
    () => {
      return {
        data: args.clusterDetails,
        error: args.error,
        isValidating: args.isValidating,
        mutate: args.mutate,
      };
    }
  );
  return (
    <DiProvider use={[useListClusterDetailsData]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};

export const LoadingCluster = TemplateLoading.bind({});
LoadingCluster.args = {
  clusterDetails: undefined,
  error: undefined,
  isValidating: true,
  mutate: jest.fn(),
} as Partial<typeof AdminClustersBox>;
LoadingCluster.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const loading = await canvas.getByTestId('loading');
  expect(loading).not.toBeNull();
};
LoadingCluster.parameters = {
  msw: {
    handlers: [],
  },
};
