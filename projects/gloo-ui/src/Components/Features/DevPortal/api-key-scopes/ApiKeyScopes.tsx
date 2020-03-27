import React from 'react';
import styled from '@emotion/styled';
import { useParams, useHistory } from 'react-router';
import { ApiKeyScopeCard } from './ApiKeyScopeCard';
import { ReactComponent as PortalIcon } from 'assets/single-portal-icon.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { DevPortalApi } from '../api';
import useSWR from 'swr';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloModal } from 'Components/Common/SoloModal';
import { EditKeyScopeModal } from './EditKeyScopeModal';
import { KeyScopeStatus } from 'proto/dev-portal/api/dev-portal/v1/portal_pb';

const APIKeyScopesBlock = styled.div`
  position: relative;
`;

const CreationButtonArea = styled.div`
  position: absolute;
  top: -42px;
  right: 0;
  cursor: pointer;
  display: flex;
  font-size: 14px;
  cursor: pointer;

  svg {
    margin-right: 8px;
    height: 22px;
    width: 22px;
  }
`;

export const APIKeyScopes = () => {
  const [keyScopeExpanded, setKeyScopeExpanded] = React.useState<{
    portalUid: string;
    scopeName: string;
  }>();
  const [createScopeWizardOpen, setCreateScopeWizardOpen] = React.useState(
    false
  );

  const { data: portalsList, error: getApiKeyDocsError } = useSWR(
    'listPortals',
    DevPortalApi.listPortals
  );

  const openCreateScope = () => {
    setCreateScopeWizardOpen(true);
  };
  const finishCreateScope = (scopeData: KeyScopeStatus.AsObject) => {
    setCreateScopeWizardOpen(false);
  };
  const cancelCreateScope = () => {
    setCreateScopeWizardOpen(false);
  };

  const expandCard = (portalUid: string, scopeName: string) => {
    setKeyScopeExpanded(
      keyScopeExpanded?.portalUid === portalUid &&
        keyScopeExpanded?.scopeName === scopeName
        ? undefined
        : {
            portalUid,
            scopeName
          }
    );
  };

  return (
    <APIKeyScopesBlock>
      <CreationButtonArea onClick={openCreateScope}>
        <span className='flex items-center text-green-400 cursor-pointer hover:text-green-300'>
          <GreenPlus className='fill-current' />
          <span className='text-gray-700'> Create a Scope</span>
        </span>
      </CreationButtonArea>
      {portalsList?.map(portalInfo => (
        <SectionCard
          key={portalInfo.metadata?.uid}
          cardName={portalInfo.metadata!.name}
          logoIcon={<PortalIcon />}>
          <div>
            {/*portalInfo.spec?.keyScopesList*/ [{ name: 'dig' }].map(
              (scopeInfo, ind) => (
                <ApiKeyScopeCard
                  key={scopeInfo.name}
                  name={scopeInfo.name}
                  onClick={() =>
                    expandCard(portalInfo.metadata!.uid, scopeInfo.name)
                  }
                  isExpanded={
                    keyScopeExpanded?.portalUid === portalInfo.metadata?.uid &&
                    keyScopeExpanded?.scopeName === scopeInfo.name
                  }
                />
              )
            )}
          </div>
        </SectionCard>
      ))}

      <SoloModal visible={createScopeWizardOpen} width={750} noPadding={true}>
        <EditKeyScopeModal
          onEdit={finishCreateScope}
          onCancel={cancelCreateScope}
          createNotEdit={true}
        />
      </SoloModal>
    </APIKeyScopesBlock>
  );
};
