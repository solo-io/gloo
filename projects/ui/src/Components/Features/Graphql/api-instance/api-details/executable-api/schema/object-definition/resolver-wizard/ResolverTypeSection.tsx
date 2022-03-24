import {
  getResolverFromConfig,
  getUpstream,
  ResolverWizardFormProps,
} from './ResolverWizard';
import {
  SoloFormDropdown,
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

const getType = (
  resolver: Resolution.AsObject
): ResolverWizardFormProps['resolverType'] => {
  if (resolver.grpcResolver) {
    return 'gRPC';
  } else if (resolver.restResolver) {
    return 'REST';
  }
  return 'REST';
};

export const ResolverTypeSection = ({ isEdit }: ResolverTypeSectionProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();
  const [selectedName, setSelectedName] = React.useState<string>();
  const resolverOptions = formik.values.listOfResolvers.map(([rName]) => {
    return {
      key: rName,
      value: rName,
    };
  });
  const onResolverCopy = (copyName: any) => {
    const { values } = formik;
    const resolver = values.listOfResolvers.find(([rName]) => {
      return rName === copyName;
    });
    if (resolver) {
      const [_rName, newResolver] = resolver;
      setSelectedName(_rName);
      const upstream = getUpstream(newResolver);
      const resolverType = getType(newResolver);
      const stringifiedResolver = getResolverFromConfig(newResolver);
      formik.setFieldValue('resolverConfig', stringifiedResolver);
      formik.setFieldValue('upstream', upstream);
      formik.setFieldValue('resolverType', resolverType);
    }
  };

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
      {formik.values.listOfResolvers.length > 0 && (
        <div
          data-testid='create-resolver-from-config'
          className='grid grid-cols-2 gap-4 '>
          <div>
            <SoloFormDropdown
              searchable={true}
              name='resolverCopy'
              title='Create Resolver From Config'
              value={selectedName}
              onChange={onResolverCopy}
              options={resolverOptions}
            />
          </div>
        </div>
      )}
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
