import React from 'react';
import { useParams } from 'react-router';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as CodeIcon } from 'assets/code-icon.svg';
import { healthConstants } from 'Styles';
import { useHistory } from 'react-router-dom';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';

export const APIDetails = () => {
  const { apiname } = useParams();
  const history = useHistory();
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Dev Portal section</div>}>
      <div>
        <SectionCard
          cardName={apiname || 'API'}
          logoIcon={
            <span className='text-blue-500'>
              <CodeIcon className='fill-current' />
            </span>
          }
          health={healthConstants.Good.value}
          headerSecondaryInformation={[
            {
              title: 'Modified',
              value: 'Feb 26, 2020'
            }
          ]}
          healthMessage={'Portal Status'}
          onClose={() => history.push(`/dev-portal/`)}></SectionCard>
      </div>
    </ErrorBoundary>
  );
};
