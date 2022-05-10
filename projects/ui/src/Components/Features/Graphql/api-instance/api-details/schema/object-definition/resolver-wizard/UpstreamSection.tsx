import { useListUpstreams } from 'API/hooks';
import { SoloFormDropdown } from 'Components/Common/SoloFormComponents';
import { useFormikContext } from 'formik';
import * as React from 'react';
import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import { getUpstreamId, getUpstreamRefFromId } from 'utils/graphql-helpers';
import YAML from 'yaml';
import { getDefaultConfigFromType } from './ResolverConfigSection';
import { ResolverWizardFormProps } from './ResolverWizard';
import * as ResolverWizardStyles from './ResolverWizard.styles';

export const UpstreamSection: React.FC<{
  existingUpstreamId: string;
  nextButtonDisabled: boolean;
  onNextClicked(): void;
  onCancel(): void;
}> = ({ existingUpstreamId, nextButtonDisabled, onNextClicked, onCancel }) => {
  const formik = useFormikContext<ResolverWizardFormProps>();
  const { data: upstreams } = useListUpstreams();

  const onChange = (newUpstreamId: any) => {
    try {
      formik.setFieldError('upstream', undefined);
      const demoConfig = getDefaultConfigFromType(formik.values.resolverType);

      if (!formik.values.resolverConfig) {
        formik.setFieldValue('resolverConfig', demoConfig);
      }

      const resolverValue = YAML.parse(
        formik.values.resolverConfig || demoConfig
      );

      // TODO: Update the form state so that the resolverConfig doesn't
      // contain the upstreamRef in the form state YAML. The form shouldn't
      // handle this logic. Add the upstreamRef in to the marshalled
      // object when the request is sent in the `graphql.ts` file.

      // Updates the upstream ref value in the YAML
      const newUpstreamRef = getUpstreamRefFromId(newUpstreamId);
      if (!!resolverValue?.restResolver?.upstreamRef)
        resolverValue.restResolver.upstreamRef = newUpstreamRef;
      if (!!resolverValue?.grpcResolver?.upstreamRef)
        resolverValue.grpcResolver.upstreamRef = newUpstreamRef;

      // nullStr makes sure that it doesn't put NULL everywhere.
      // simpleKeys makes sure that `:method:` doesn't become `? method:`.
      YAML.scalarOptions.null.nullStr = '';
      const stringifiedResolver = YAML.stringify(resolverValue, {
        simpleKeys: true,
      });
      formik.setFieldValue('resolverConfig', stringifiedResolver);
    } catch (error) {
      console.error({ error });
    }
  };

  return (
    <>
      <div data-testid='upstream-section' className='w-full h-full px-6 pb-0'>
        <div className='grid gap-4 '>
          <div className='mb-2 '>
            <label className='text-base font-medium '>Upstream</label>
            <div className='mt-3'>
              <SoloFormDropdown
                data-testid='upstream-section-dropdown'
                name='upstream'
                defaultValue={existingUpstreamId}
                searchable={true}
                onChange={onChange}
                options={upstreams
                  ?.map(upstream => {
                    return {
                      key: upstream.metadata?.uid!,
                      value: getUpstreamId(upstream.metadata),
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

      <div className='flex items-center justify-between px-6'>
        <ResolverWizardStyles.IconButton onClick={onCancel}>
          Cancel
        </ResolverWizardStyles.IconButton>
        <SoloButtonStyledComponent
          onClick={onNextClicked}
          disabled={nextButtonDisabled}>
          Next Step
        </SoloButtonStyledComponent>
      </div>
    </>
  );
};
