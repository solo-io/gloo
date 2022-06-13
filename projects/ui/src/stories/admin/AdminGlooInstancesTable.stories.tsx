import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { AdminGlooInstancesTable } from '../../Components/Features/Admin/AdminGlooInstancesTable';
import { MemoryRouter } from 'react-router';
import {
  createGateway,
  createGlooInstance,
  createGlooInstanceObj,
} from 'stories/mocks/generators';
import { expect, jest } from '@storybook/jest';
import * as Apis from 'API/hooks';
import { within } from '@storybook/testing-library';
import { DiProvider, injectable } from 'react-magnetic-di/macro';

export default {
  title: `Admin / ${AdminGlooInstancesTable.name}`,
  component: AdminGlooInstancesTable,
} as unknown as ComponentMeta<typeof AdminGlooInstancesTable>;

const Template: ComponentStory<typeof AdminGlooInstancesTable> = args => {
  const useListGlooInstancesDi = injectable(Apis.useListGlooInstances, () => {
    const data = Array.from({ length: 1 }).map(() => {
      return createGlooInstanceObj();
    });
    return { data, error: undefined, isValidating: false, mutate: jest.fn() };
  });

  const useListGatewaysDi = injectable(Apis.useListGateways, () => {
    const data = Array.from({ length: 1 }).map(() => {
      return createGateway();
    });
    return { data, error: undefined, isValidating: false, mutate: jest.fn() };
  });

  return (
    <MemoryRouter>
      <DiProvider use={[useListGlooInstancesDi, useListGatewaysDi]}>
        <AdminGlooInstancesTable />
      </DiProvider>
    </MemoryRouter>
  );
};

export const Primary = Template.bind({});
Primary.args = {} as Partial<typeof AdminGlooInstancesTable>;
Primary.play = async ({ canvasElement }) => {
  const canvas = within(canvasElement);
  const summary = await canvas.getByTestId('admin-gloo-instances-table');
  expect(summary).not.toBeNull();
};
