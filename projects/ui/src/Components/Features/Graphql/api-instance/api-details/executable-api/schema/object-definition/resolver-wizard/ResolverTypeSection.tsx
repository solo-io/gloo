import { ResolverWizardFormProps } from './ResolverWizard';
import {
  SoloFormRadio,
  SoloFormRadioOption,
} from 'Components/Common/SoloFormComponents';
import React from 'react';
import { useFormikContext } from 'formik';
import { Resolution } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';

export type ResolverTypeSectionProps = { isEdit: boolean };

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
] as SoloFormRadioOption[];

export const getType = (
  resolver?: Resolution.AsObject
): ResolverWizardFormProps['resolverType'] => {
  if (!resolver) {
    return 'REST';
  }
  if (resolver.grpcResolver) {
    return 'gRPC';
  } else if (resolver.restResolver) {
    return 'REST';
  }
  return 'REST';
};

export const ResolverTypeSection = ({ isEdit }: ResolverTypeSectionProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();

  const onTypeChange = (resolverType: string) => {
    formik.setFieldValue('resolverType', resolverType);
    formik.setFieldValue('resolverConfig', '');
    formik.setFieldValue('upstream', '');
  };

  return (
    <div data-testid='resolver-type-section' className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>

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
