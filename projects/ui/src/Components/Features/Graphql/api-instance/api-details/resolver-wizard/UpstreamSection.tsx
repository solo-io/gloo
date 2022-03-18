import * as React from 'react';
import { ResolverWizardFormProps } from '../ResolverWizard';
import { useFormikContext } from 'formik';
import { useListUpstreams } from 'API/hooks';
import { SoloFormDropdown } from 'Components/Common/SoloFormComponents';
import YAML from 'yaml';
import { getDefaultConfigFromType } from './ResolverConfigSection';

type UpstreamSectionProps = {
  isEdit: boolean;
  existingUpstream?: string;
};

export const UpstreamSection = ({
  isEdit,
  existingUpstream,
}: UpstreamSectionProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();
  const { data: upstreams } = useListUpstreams();

  const onChange = (value: any) => {
    formik.setFieldValue('upstream', value);
    try {
      const demoConfig = getDefaultConfigFromType(formik.values.resolverType);
      if (!formik.values.resolverConfig) {
        formik.setFieldValue('resolverConfig', demoConfig);
      }
      const resolverValue = YAML.parse(formik.values.resolverConfig || demoConfig);
      if (resolverValue?.restResolver?.upstreamRef) {
        const [name, namespace] = value.split('::');
        resolverValue.restResolver.upstreamRef.name = name;
        resolverValue.restResolver.upstreamRef.namespace = namespace;
        const stringifiedResolver = YAML.stringify(resolverValue);
        formik.setFieldValue('resolverConfig', stringifiedResolver);
      }
    } catch (ignore) { }
  }

  return (
    <div data-testid='upstream-section' className='w-full h-full p-6 pb-0'>
      <div
        className={'flex items-center mb-6 text-lg font-medium text-gray-800'}>
        {isEdit ? 'Edit' : 'Configure'} Resolver{' '}
      </div>
      <div className='grid gap-4 '>
        <div className='mb-2 '>
          <label className='text-base font-medium '>Upstream</label>
          <div className='mt-3'>
            <SoloFormDropdown
              name='upstream'
              value={formik.values.upstream}
              defaultValue={existingUpstream}
              searchable={true}
              onChange={onChange}
              options={upstreams
                ?.map(upstream => {
                  return {
                    key: upstream.metadata?.uid!,
                    value: `${upstream.metadata?.name!}::${
                      upstream.metadata?.namespace
                    }`,
                    displayValue: upstream.metadata?.name!,
                  };
                })
                .sort((upstream1, upstream2) =>
                  upstream1.displayValue === upstream2.displayValue
                    ? 0
                    : (upstream1?.displayValue ?? upstream1.value) >
                      (upstream2?.displayValue ?? upstream2.value)
                    ? 1
                    : -1
                )}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
