import styled from '@emotion/styled';
import { ReactComponent as PortalIcon } from 'assets/single-portal-icon.svg';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { SoloModal } from 'Components/Common/SoloModal';
import { KeyScopeStatus } from '@solo-io/dev-portal-grpc/dev-portal/api/dev-portal/v1/portal_pb';
import { ApiKeyScopeWithApiDocs } from '@solo-io/dev-portal-grpc/dev-portal/api/grpc/admin/api_key_scope_pb';
import React from 'react';
import useSWR from 'swr';
import { apiKeyScopeApi } from '../api';
import { NoDataPanel } from '../DevPortal';
import { ApiKeyScopeCard } from './ApiKeyScopeCard';
import { EditKeyScopeModal } from './EditKeyScopeModal';

const CardContainer = styled.div`
  margin-bottom: 12px;
`;

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
    portalRef: string;
    scopeName: string;
  }>();
  const [createScopeWizardOpen, setCreateScopeWizardOpen] = React.useState(
    false
  );

  const { data: keyScopesList, error: listKeyScopesError } = useSWR(
    'listKeyScopes',
    apiKeyScopeApi.listKeyScopes
  );

  type PortalToKeyScope = {
    [portalRef: string]: ApiKeyScopeWithApiDocs.AsObject[];
  };

  const portalToKeyScope: PortalToKeyScope = (keyScopesList || []).reduce(
    (acc: PortalToKeyScope, keyScopeResponse) => {
      const portalRef = `${keyScopeResponse.apiKeyScope!.portal!.namespace}.${
        keyScopeResponse.apiKeyScope!.portal!.name
      }`;

      let v = {
        ...acc,
        [portalRef]: [
          keyScopeResponse,
          ...(!!acc[portalRef] ? acc[portalRef] : [])
        ]
      };
      return v;
    },
    {}
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

  const expandCard = (portalRef: string, scopeName: string) => {
    setKeyScopeExpanded(
      keyScopeExpanded?.portalRef === portalRef &&
        keyScopeExpanded?.scopeName === scopeName
        ? undefined
        : {
            portalRef,
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
      {!Object.keys(portalToKeyScope).length && (
        <NoDataPanel
          missingContentText='There are no API Key Scopes to display'
          helpText='Create a Key Scope to allow users to generate API Keys.'
          identifier='key-scopes-page'
        />
      )}
      {!!Object.keys(portalToKeyScope).length &&
        Object.keys(portalToKeyScope)
          ?.sort((a, b) => (a === b ? 0 : a > b ? 1 : -1))
          .map(portalRef => (
            <SectionCard
              key={portalRef}
              // TODO joekelley get the portal display name
              cardName={portalRef.split('.')[1]}
              logoIcon={<PortalIcon />}>
              <div>
                {portalToKeyScope[portalRef]
                  .sort((a, b) =>
                    a.apiKeyScope!.spec!.name === b.apiKeyScope!.spec!.name
                      ? 0
                      : a.apiKeyScope!.spec!.name > b.apiKeyScope!.spec!.name
                      ? 1
                      : -1
                  )
                  .map(scopeInfo => (
                    <CardContainer key={scopeInfo.apiKeyScope?.status?.name}>
                      <ApiKeyScopeCard
                        onClick={() =>
                          expandCard(
                            portalRef,
                            scopeInfo.apiKeyScope!.status!.name
                          )
                        }
                        isExpanded={
                          keyScopeExpanded?.portalRef === portalRef &&
                          keyScopeExpanded?.scopeName ===
                            scopeInfo.apiKeyScope?.status?.name
                        }
                        keyScope={scopeInfo}
                      />
                    </CardContainer>
                  ))}
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
