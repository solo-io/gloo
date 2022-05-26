import {
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
  usePageApiRef,
} from 'API/hooks';
import { ReactComponent as RouteIcon } from 'assets/route-icon.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import React, { useState } from 'react';
import { Spacer } from 'Styles/StyledComponents/spacer';
import { ResolverWizard } from './resolver-wizard/ResolverWizard';

const ResolverButton: React.FC<{
  field: any;
  objectType: string;
  resolveDirectiveExists: boolean;
}> = ({ field, objectType, resolveDirectiveExists }) => {
  const apiRef = usePageApiRef();
  const { mutate: mutateDetails } = useGetGraphqlApiDetails(apiRef);
  const { mutate: mutateYaml } = useGetGraphqlApiYaml(apiRef);
  const [showResolverWizard, setShowResolverWizard] = useState(false);
  return (
    <div className='text-left w-100'>
      <Spacer
        padding={1}
        px={2}
        data-testid={`resolver-${field?.name?.value}`}
        className={`inline-flex items-center ${
          resolveDirectiveExists
            ? 'focus:ring-blue-500gloo text-blue-700gloo bg-blue-200gloo  border-blue-600gloo hover:bg-blue-300gloo'
            : 'focus:ring-gray-500 text-gray-700 bg-gray-300  border-gray-600 hover:bg-gray-200'
        }
            border rounded-full shadow-sm cursor-pointer focus:outline-none focus:ring-2 focus:ring-offset-2 `}
        onClick={() => setShowResolverWizard(true)}>
        <RouteIcon
          data-testid={`route-${field?.name?.value}`}
          className='w-5 h-5 mr-1 fill-current text-blue-600gloo'
        />
        {resolveDirectiveExists ? 'Edit Resolver' : 'Define Resolver'}
      </Spacer>

      {showResolverWizard && (
        <SoloModal visible={showResolverWizard} width={750}>
          <ResolverWizard
            apiRef={apiRef}
            field={field}
            objectType={objectType}
            onClose={() => {
              setTimeout(() => {
                mutateDetails();
                mutateYaml();
              }, 300);
              setShowResolverWizard(false);
            }}
          />
        </SoloModal>
      )}
    </div>
  );
};

export default ResolverButton;
