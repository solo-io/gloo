import { ResolverWizardFormProps } from '../ResolverWizard';
import {
  SoloFormRadio,
  SoloFormRadioOption,
} from 'Components/Common/SoloFormComponents';
import React from 'react';

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

export const ResolverTypeSection = ({ isEdit }: ResolverTypeSectionProps) => {
  return (
    <div data-testid='resolver-type-section' className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className='grid grid-cols-2 gap-4 '>
        <SoloFormRadio<ResolverWizardFormProps>
          name='resolverType'
          isUpdate={Boolean(isEdit)}
          title='Resolver Type'
          options={apiTypeOptions}
          titleAbove
        />
      </div>
    </div>
  );
};
