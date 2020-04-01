import * as React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles';
import { ReactComponent as StacksIcon } from 'assets/app-icon.svg';
import { ReactComponent as CodeGearIcon } from 'assets/code-sprocket-icon.svg';
import { ReactComponent as EyeIcon } from 'assets/view-icon.svg';
import useSWR from 'swr';
import { apiKeyScopeApi } from '../api';
import { SoloToggleSwitch } from 'Components/Common/SoloToggleSwitch';
import {
  SoloNegativeButton,
  SoloButtonStyledComponent
} from 'Styles/CommonEmotions/button';
//import { SoloModal } from './SoloModal';
import { ConfirmationModal } from 'Components/Common/ConfirmationModal';
import { SoloModal } from 'Components/Common/SoloModal';
import { EditKeyScopeModal } from './EditKeyScopeModal';
import { ApiDoc } from 'proto/dev-portal/api/grpc/admin/apidoc_pb';
import { ApiKeyScopeWithApiDocs } from 'proto/dev-portal/api/grpc/admin/api_key_scope_pb';

const Card = styled.div`
  background: white;
  border-radius: 10px;
  box-shadow: 0 0 10px ${colors.darkerBoxShadow};
`;

type ExpandedProps = {
  isExpanded: boolean;
};

const KeyCoreInformation = styled.div`
  display: grid;
  grid-template-columns: 1fr 6fr 1fr 4fr;
  padding: 25px 20px 20px;
  cursor: pointer;
  ${(props: ExpandedProps) =>
    props.isExpanded
      ? `border-radius: 10px 10px 0 0;`
      : 'border-radius: 10px;'};
`;

const StacksImageHolder = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 83px;
  height: 83px;
  border-radius: 83px;
  background: ${colors.februaryGrey};
  margin: auto 30px auto 10px;

  svg {
    width: 55px;
  }
`;

const AppTitleArea = styled.div`
  font-size: 16px;
  line-height: 19px;
`;
const AppTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
  margin-bottom: 5px;
  font-weight: 500;
`;
const AppDescription = styled.div``;
const AppCreatedOnDate = styled.div`
  color: ${colors.juneGrey};
`;
const Divider = styled.div`
  width: 1px;
  height: 100%;
  background: ${colors.marchGrey};
  margin: 0 85px;
`;
const CountArea = styled.div`
  display: flex;
  align-items: center;
  margin-right: 70px;
  font-size: 16px;
  line-height: 19px;
`;
const CountIconCircle = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 42px;
  height: 42px;
  border-radius: 42px;
  background: ${colors.seaBlue};

  svg * {
    fill: white;
  }
`;
const CountNumeral = styled.span`
  font-weight: bold;
  margin: 0 3px 0 8px;
`;

const KeyExtraInformation = styled.div`
  position: relative;
  background: ${colors.januaryGrey};
  border-top: 1px solid ${colors.marchGrey};
  border-radius: 0 0 10px 10px;
  overflow: hidden;
  ${(props: ExpandedProps) =>
    props.isExpanded
      ? `max-height: 1000px; transition: max-height .7s ease-in;`
      : `max-height: 0; transition: max-height .7s ease-out;`};

  > div {
    padding: 20px 20px 25px;
  }
`;
const ButtonActionGroup = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
  border-top: 1px solid ${colors.marchGrey};
  background: white;
  border-radius: 0 0 10px 10px;

  > button {
    margin-left: 8px;
  }
`;

const NothingFoundInfo = styled.div`
  font-size: 16px;
  padding: 8px 20px;
  background: white;
  border: 1px solid ${colors.marchGrey};
  border-radius: 8px;
`;

const AccessList = styled.div``;
const AccessListTitle = styled.div`
  font-size: 22px;
  line-height: 26px;
  margin-bottom: 15px;
  margin-left: 3px;
`;
const APIAccessBlock = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  line-height: 19px;
  height: 40px;
  padding: 0 20px;
  background: white;
  border: 1px solid ${colors.marchGrey};
  border-radius: 8px;
`;
const AccessIconHolder = styled.div`
  display: flex;
  align-items: center;
  margin-right: 5px;

  svg {
    height: 20px;
    fill: ${colors.seaBlue};
  }
`;
const AccessName = styled.div`
  font-weight: 500;
`;
const AccessVersion = styled.div`
  font-size: 12px;
  line-height: 12px;
  font-weight: 500;
  padding: 2px 4px;
  border: 1px solid ${colors.juneGrey};
  border-radius: 12px;
  margin-left: 8px;
`;
const AccessDescription = styled.div`
  flex: 1;
  color: ${colors.juneGrey};
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 700px;
`;
const AccessToggleHolder = styled.div``;

const EmptyContentBlock = styled.div`
  display: flex;
  padding: 19px 20px;
`;

const EmptyDescriptorsArea = styled.div`
  width: 350px;
`;

const EmptyTitleLine = styled.div`
  width: 175px;
  height: 12px;
  background: ${colors.septemberGrey};
  margin-bottom: 18px;
`;
const EmptyDescriptionLine = styled.div`
  max-width: 350px;
  height: 12px;
  background: ${colors.aprilGrey};
  margin-bottom: 14px;
`;
const EmptyCountArea = styled.div`
  display: flex;
  align-items: center;
  width: 175px;

  div:nth-child(2) {
    margin: 0;
    margin-left: 12px;
    width: 100px;
  }
`;

interface APICardProps {
  keyScope: ApiKeyScopeWithApiDocs.AsObject;
  onClick?: () => any;
  isExpanded?: boolean;
}

export const ApiKeyScopeCard = ({
  keyScope,
  onClick,
  isExpanded
}: APICardProps) => {
  const [editWizardOpen, setEditWizardOpen] = React.useState(false);
  const [attemptingDelete, setAttemptingDelete] = React.useState(false);

  const beginEditing = () => {
    setEditWizardOpen(true);
  };
  const cancelEditing = () => {
    setEditWizardOpen(false);
  };
  const finishEditing = () => {
    setEditWizardOpen(false);
  };

  const attemptDeleteScope = () => {
    setAttemptingDelete(true);
  };

  const deleteScope = () => {
    apiKeyScopeApi.deleteKeyScope({
      portal: {
        namespace: keyScope.apiKeyScope!.portal!.namespace,
        name: keyScope.apiKeyScope!.portal!.name
      },
      apiKeyScope: {
        namespace: keyScope.apiKeyScope!.spec!.namespace,
        name: keyScope.apiKeyScope!.spec!.name
      }
    });
    setAttemptingDelete(false);
  };

  const cancelDeletion = () => {
    setAttemptingDelete(false);
  };

  return (
    <>
      <Card>
        <KeyCoreInformation isExpanded={!!isExpanded} onClick={onClick}>
          <StacksImageHolder>
            <StacksIcon />
          </StacksImageHolder>
          <AppTitleArea>
            <AppTitle>
              {keyScope.apiKeyScope?.spec?.displayName ||
                keyScope.apiKeyScope?.spec?.name}
            </AppTitle>
            <AppDescription>
              {keyScope.apiKeyScope?.spec?.description}
            </AppDescription>
          </AppTitleArea>
          <Divider />
          <CountArea>
            <CountIconCircle>
              <CodeGearIcon />
            </CountIconCircle>
            <CountNumeral>{keyScope.apiDocsList.length}</CountNumeral> APIs in{' '}
            scope
          </CountArea>
        </KeyCoreInformation>
        <KeyExtraInformation isExpanded={!!isExpanded}>
          <AccessList>
            <AccessListTitle>API Access</AccessListTitle>
            {!!keyScope.apiDocsList.length ? (
              keyScope.apiDocsList.map(apiDoc => (
                <APIAccessBlock key={apiDoc.metadata?.uid}>
                  <AccessIconHolder>
                    <CodeGearIcon />
                  </AccessIconHolder>
                  <AccessName>{apiDoc.status?.displayName}</AccessName>
                  <AccessVersion>{apiDoc.status?.version}</AccessVersion>
                  <AccessDescription>
                    {apiDoc.status?.description}
                  </AccessDescription>
                </APIAccessBlock>
              ))
            ) : (
              <NothingFoundInfo>
                No access has been granted for this scope.
              </NothingFoundInfo>
            )}
          </AccessList>
          <ButtonActionGroup>
            <SoloButtonStyledComponent onClick={beginEditing}>
              Edit Scope
            </SoloButtonStyledComponent>
            <SoloNegativeButton onClick={attemptDeleteScope}>
              Delete Scope
            </SoloNegativeButton>
          </ButtonActionGroup>
        </KeyExtraInformation>
      </Card>

      <SoloModal visible={editWizardOpen} width={750} noPadding={true}>
        <EditKeyScopeModal
          existingKeyScope={keyScope.apiKeyScope}
          onEdit={finishEditing}
          onCancel={cancelEditing}
        />
      </SoloModal>
      <ConfirmationModal
        visible={attemptingDelete}
        confirmationTopic='Delete this scope?'
        confirmText='Delete'
        goForIt={deleteScope}
        cancel={cancelDeletion}
        isNegative={true}
      />
    </>
  );
};

export const EmptyKeyScopeCard = () => (
  <EmptyContentBlock>
    <>
      <StacksImageHolder>
        <StacksIcon />
      </StacksImageHolder>
      <EmptyDescriptorsArea>
        <EmptyTitleLine />
        <EmptyDescriptionLine />
        <div style={{ width: '175px' }}>
          <EmptyDescriptionLine />
        </div>
      </EmptyDescriptorsArea>
      <Divider />

      <EmptyCountArea>
        <CountIconCircle>
          <CodeGearIcon />
        </CountIconCircle>
        <EmptyDescriptionLine />
      </EmptyCountArea>
    </>
  </EmptyContentBlock>
);
