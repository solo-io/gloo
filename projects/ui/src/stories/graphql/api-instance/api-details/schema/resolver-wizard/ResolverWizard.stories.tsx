import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import {
  ResolverWizard,
  ResolverWizardProps,
} from '../../../../../../Components/Features/Graphql/api-instance/api-details/schema/resolver-wizard/ResolverWizard';
import { Formik } from 'formik';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { mockEnumDefinitions } from 'Components/Features/Graphql/api-instance/api-details/schema/mockData';
import { expect, jest } from '@storybook/jest';
import { DiProvider, injectable } from 'react-magnetic-di/macro';
import {
  useConfirm,
  ConfirmOptions,
} from 'Components/Context/ConfirmModalContext';
import * as Apis from 'API/hooks';
import { useAppSettings } from 'Components/Context/AppSettingsContext';
import { GraphqlApi } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import {
  createGlooInstance,
  createGlooInstanceObj,
  createInstanceStatus,
  createObjMeta,
  createObjRef,
  createObjSpec,
} from 'stories/mocks/generators';

type StoryResolverWizard = ResolverWizardProps & {
  apiDetails: GraphqlApi.AsObject;
  consoleOptions: {
    readonly: boolean;
    apiExplorerEnabled: boolean;
    errorMessage: string;
  };
};

export default {
  title:
    'Graphql / api-instance / api-details / schema / resolver-wizard / ResolverWizard',
  component: ResolverWizard,
} as ComponentMeta<typeof ResolverWizard>;

const Template: ComponentStory<typeof ResolverWizard> = (
  args: StoryResolverWizard | any
) => {
  // GraphqlApi.AsObject
  const newArgs: typeof ResolverWizard & StoryResolverWizard = { ...args };
  const useGetGraphqlApiDetailsDi = injectable(
    Apis.useGetGraphqlApiDetails,
    () => {
      return {
        data: newArgs.apiDetails,
        isValidating: false,
        mutate: jest.fn(),
        error: undefined,
      };
    }
  );
  const useGetConsoleOptionsDi = injectable(Apis.useGetConsoleOptions, () => {
    return {
      readonly: false,
      apiExplorerEnabled: true,
      errorMessage: '',
    };
  });
  const useConfirmDi = injectable(useConfirm, () => {
    return jest.fn((options: ConfirmOptions) => {
      return new Promise(() => {
        return jest.fn();
      });
    });
  });

  const useAppSettingsDi = injectable(useAppSettings, () => {
    return {
      appSettings: {
        keyboardHandler: '',
      },
      onAppSettingsChange: jest.fn(),
    };
  });

  return (
    <DiProvider
      use={[
        useGetGraphqlApiDetailsDi,
        useGetConsoleOptionsDi,
        useConfirmDi,
        useAppSettingsDi,
      ]}>
      <Formik
        onSubmit={jest.fn()}
        initialValues={{
          resolverType: {},
          upstream: '',
          resolverConfig: '',
          listOfResolvers: [],
          protoFile: '',
        }}>
        <ResolverWizard {...args} />
      </Formik>
    </DiProvider>
  );
};

export const Primary = Template.bind({});

const apiRef = new ClusterObjectRef();

const field = { ...mockEnumDefinitions[0] };
field.directives = [];
(field as any).definitions = [];

const metadata = createObjMeta();
const glooInstance = createObjRef();

Primary.args = {
  apiRef: apiRef.toObject(),
  apiDetails: {
    metadata,
    spec: undefined,
    status: undefined,
    glooInstance,
  },
  field,
  objectType: '',
} as Partial<typeof ResolverWizard>;
