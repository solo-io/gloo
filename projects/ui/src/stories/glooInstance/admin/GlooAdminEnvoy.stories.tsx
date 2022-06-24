import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { GlooAdminEnvoy } from '../../../Components/Features/GlooInstance/Admin/GlooAdminEnvoy';
import * as Apis from 'API/hooks';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import { jest } from '@storybook/jest';
import { faker } from '@faker-js/faker';

import { createConfigDump } from 'stories/mocks/generators';
import { useParams } from 'react-router';

export default {
  title: 'GlooInstance / Admin / GlooAdminEnvoy',
  component: GlooAdminEnvoy,
} as ComponentMeta<typeof GlooAdminEnvoy>;

const Template: ComponentStory<typeof GlooAdminEnvoy> = (args: any) => {
  const useParamsDi = injectable(useParams, () => {
    return {
      name: args.name,
      namespace: args.namespace,
    };
  });
  const useGetConfigDumpsDi = injectable(Apis.useGetConfigDumps, () => {
    return {
      data: args.data,
      error: args.error,
      mutate: args.mutate,
      isValidating: args.isValidating,
    };
  });
  return (
    <DiProvider use={[useParamsDi, useGetConfigDumpsDi]}>
      <GlooAdminEnvoy />
    </DiProvider>
  );
};

const data = Array.from({ length: 1 }).map(() => {
  return createConfigDump({ error: '' });
});

export const Primary = Template.bind({});
Primary.args = {
  data,
  name: faker.random.word(),
  namespace: faker.random.word(),
  mutate: jest.fn(),
  isValidating: false,
  error: undefined,
} as Partial<typeof GlooAdminEnvoy>;
