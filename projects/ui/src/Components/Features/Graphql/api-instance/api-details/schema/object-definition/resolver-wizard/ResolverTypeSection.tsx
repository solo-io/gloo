import { ResolverWizardFormProps } from './ResolverWizard';
import {
  SoloFormRadio,
  SoloFormRadioOption,
} from 'Components/Common/SoloFormComponents';
import React from 'react';
import { useFormikContext } from 'formik';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import { getDefaultConfigFromType } from './ResolverConfigSection';

export let apiTypeOptions = [
  {
    displayValue: 'REST',
    value: 'REST',
    subHeader:
      'Integrate with upstream REST APIs and customize HTTP request and response mappings.',
  },
  {
    displayValue: 'gRPC',
    value: 'gRPC',
    subHeader: 'Integrate with upstream gRPC APIs based on a proto definition.',
  },
  {
    displayValue: 'Mock',
    value: 'Mock',
    subHeader:
      'Create a mock response value for the GraphlApi instead of wiring it to an upstream API. This can be useful for demos and rapid prototyping.',
  },
] as SoloFormRadioOption[];

export const getType = (
  resolver?: Resolution.AsObject
): ResolverWizardFormProps['resolverType'] => {
  if (resolver?.mockResolver) return 'Mock';
  if (resolver?.grpcResolver) return 'gRPC';
  if (resolver?.restResolver) return 'REST';
  return 'REST';
};

export const ResolverTypeSection: React.FC<{
  setWarningMessage(message: string): void;
}> = ({ setWarningMessage }) => {
  const formik = useFormikContext<ResolverWizardFormProps>();

  const onTypeChange = (
    resolverType: ResolverWizardFormProps['resolverType']
  ) => {
    setWarningMessage('');
    formik.setFieldValue('resolverType', resolverType);
    if (resolverType !== 'gRPC') {
      formik.setFieldValue('protoFile', '');
    }
    formik.setFieldValue(
      'resolverConfig',
      getDefaultConfigFromType(resolverType)
    );
    formik.setFieldValue('upstream', '');
  };

  return (
    <div
      data-testid='resolver-type-section'
      className='w-full h-full px-6 pb-0'>
      <div className='grid grid-cols-2 gap-4 '>
        <SoloFormRadio<ResolverWizardFormProps>
          name='resolverType'
          title='Resolver Type'
          options={apiTypeOptions}
          onChange={onTypeChange}
          titleAbove
        />
      </div>
    </div>
  );
};
