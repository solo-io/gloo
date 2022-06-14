import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminGlooInstancesTable } from '../../Components/Features/Admin/AdminGlooInstancesTable';
import { MemoryRouter } from 'react-router';
import { createGateway, createGlooInstanceObj } from 'stories/mocks/generators';
import { expect, jest } from '@storybook/jest';
import * as Apis from 'API/hooks';
import { within } from '@storybook/testing-library';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { Gateway } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';

export default {
  title: 'Admin / AdminGlooInstancesTable',
  component: AdminGlooInstancesTable,
} as ComponentMeta<typeof AdminGlooInstancesTable>;

type AdminGlooInstancesTableType = {
  listGlooInstancesData: GlooInstance.AsObject[];
  listGatewaysData: Gateway.AsObject[];
};

const Template: ComponentStory<typeof AdminGlooInstancesTable> = (
  args: AdminGlooInstancesTableType | any
) => {
  const useListGlooInstancesDi = injectable(Apis.useListGlooInstances, () => {
    return {
      data: args.listGlooInstancesData,
      error: undefined,
      isValidating: false,
      mutate: jest.fn(),
    };
  });

  const useListGatewaysDi = injectable(Apis.useListGateways, () => {
    return {
      data: args.listGatewaysData,
      error: undefined,
      isValidating: false,
      mutate: jest.fn(),
    };
  });

  return (
    <MemoryRouter>
      <DiProvider use={[useListGlooInstancesDi, useListGatewaysDi]}>
        <AdminGlooInstancesTable />
      </DiProvider>
    </MemoryRouter>
  );
};

const listGlooInstancesData = Array.from({ length: 1 }).map(() => {
  return createGlooInstanceObj();
});

const listGatewaysData = Array.from({ length: 1 }).map(() => {
  return createGateway();
});

export const Primary = Template.bind({});
Primary.args = {
  listGlooInstancesData,
  listGatewaysData,
} as Partial<typeof AdminGlooInstancesTable>;
Primary.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-gloo-instances-table');
  expect(summary).not.toBeNull();
};
