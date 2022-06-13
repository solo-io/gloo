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

export default {
  title: `Admin / ${AdminClustersBox.name}`,
  component: AdminClustersBox,
} as unknown as ComponentMeta<typeof AdminClustersBox>;

const Template: ComponentStory<typeof AdminClustersBox> = args => {
  const useQueryDi = injectable(Apis.useListClusterDetails, () => {
    const details = Array.from({ length: 1 }).map(() => {
      return createClusterDetailsObj();
    });
    return {
      data: details,
      error: undefined,
      mutate: jest.fn(),
      isValidating: false,
    };
  });
  return (
    <DiProvider use={[useQueryDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};

export const OneCluster = Template.bind({});
OneCluster.args = {} as Partial<typeof AdminClustersBox>;
OneCluster.parameters = {};
OneCluster.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-box-summary');
  expect(summary).not.toBeNull();
};

const TemplateTwo: ComponentStory<typeof AdminClustersBox> = args => {
  const useQueryDi = injectable(Apis.useListClusterDetails, () => {
    const details = Array.from({ length: 2 }).map(() => {
      return createClusterDetailsObj();
    });
    return {
      data: details,
      error: undefined,
      mutate: jest.fn(),
      isValidating: false,
    };
  });
  return (
    <DiProvider use={[useQueryDi]}>
      <MemoryRouter>
        {/* @ts-ignore */}
        <AdminClustersBox {...args} />
      </MemoryRouter>
    </DiProvider>
  );
};
export const TwoCluster = TemplateTwo.bind({});

TwoCluster.args = {} as Partial<typeof AdminClustersBox>;
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

const TemplateError: ComponentStory<typeof AdminClustersBox> = args => {
  const useQueryDi = injectable(Apis.useListClusterDetails, () => {
    const meta = new grpc.Metadata();
    meta.append('thing-happened', 'A bad thing happened')!;
    const error: ServiceError = {
      message: 'Some Error Happened!',
      code: 1,
      metadata: meta,
    };
    return { data: [], error, mutate: jest.fn(), isValidating: false };
  });
  return (
    <DiProvider use={[useQueryDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};
export const ErrorCluster = TemplateError.bind({});
ErrorCluster.args = {} as Partial<typeof AdminClustersBox>;
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

const TemplateLoading: ComponentStory<typeof AdminClustersBox> = args => {
  const useQueryDi = injectable(Apis.useListClusterDetails, () => {
    return {
      data: undefined,
      error: undefined,
      isValidating: true,
      mutate: jest.fn(),
    };
  });
  return (
    <DiProvider use={[useQueryDi]}>
      <MemoryRouter>
        <AdminClustersBox />
      </MemoryRouter>
    </DiProvider>
  );
};

export const LoadingCluster = TemplateLoading.bind({});
LoadingCluster.args = {
  handler: {},
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
