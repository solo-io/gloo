import { css } from '@emotion/core';
import styled from '@emotion/styled';
import { ReactComponent as SecretsIcon } from 'assets/key-on-ring.svg';
import { ReactComponent as SettingsIcon } from 'assets/settings-gear.svg';
import { GoodStateCongratulations } from 'Components/Common/DisplayOnly/GoodStateCongratulations';
import { StatusTile } from 'Components/Common/DisplayOnly/StatusTile';
import { TallyInformationDisplay } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import React from 'react';
import { configAPI } from 'store/config/api';
import { secretAPI } from 'store/secrets/api';
import { CardCSS } from 'Styles/CommonEmotions/card';
import useSWR, { mutate } from 'swr';
import { ErrorBoundary } from '../Errors/ErrorBoundary';

const Container = styled.div`
  ${CardCSS};
`;

export const SettingsOverview = () => {
  const { data: settingsDetails, error: settingsError } = useSWR(
    'getSettings',
    configAPI.getSettings
  );

  const { data: secretsList, error: secretsError } = useSWR(
    'listSecrets',
    secretAPI.getSecretsList
  );

  const { data: watchedNamespacesList, error: watchedNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  if (!settingsDetails || !secretsList || !watchedNamespacesList) {
    return <div>Loading...</div>;
  }

  function handleSaveYAML(editedYaml: string) {
    mutate(
      'getSettings',
      configAPI.updateSettingsYaml({
        editedYamlData: {
          ref: {
            name: settingsDetails?.settings?.metadata?.name!,
            namespace: settingsDetails?.settings?.metadata?.namespace!
          },
          editedYaml
        }
      })
    );
  }
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Settings section</div>}>
      <Container>
        <div>
          <div
            css={css`
              font-size: 22px;
              margin-bottom: 10px;
            `}>
            Gloo Configuration
          </div>
          <div
            css={css`
              display: grid;
              grid-template-columns: minmax(200px, 1fr) minmax(200px, 1fr);
              grid-gap: 20px;
            `}>
            <div>
              <StatusTile
                titleText={'Settings'}
                titleIcon={
                  <SettingsIcon
                    css={css`
                      .settings-gear-a {
                        fill: #54b7e3;
                        stroke: #54b7e3;
                      }
                    `}
                  />
                }
                description={`Represents global settings for all of Gloo's components. `}
                exploreMoreLink={{
                  prompt: 'View Setttings',
                  link: `/settings/${settingsDetails.settings?.metadata?.namespace}/${settingsDetails.settings?.metadata?.name}`
                }}
                healthStatus={1}
                descriptionMinHeight={'65px'}>
                <>
                  {settingsDetails?.settings?.status?.state === 2 ? (
                    <TallyInformationDisplay
                      tallyCount={1}
                      tallyDescription={`settings error`}
                      color='orange'
                      moreInfoLink={{
                        prompt: 'View',
                        link: `/settings/${settingsDetails.settings?.metadata?.namespace}/${settingsDetails.settings?.metadata?.name}`
                      }}
                    />
                  ) : (
                    <GoodStateCongratulations typeOfItem={'settings'} />
                  )}
                  <TallyInformationDisplay
                    tallyCount={settingsDetails?.settings?.status?.state}
                    tallyDescription={`settings configuration`}
                    color='blue'
                  />
                </>
              </StatusTile>
            </div>
            <div>
              <StatusTile
                titleText={'Secrets'}
                titleIcon={
                  <span className='text-blue-600'>
                    <SecretsIcon className='fill-current' />
                  </span>
                }
                description={
                  'Certain features such as the AWS Lambda option require the use of secrets for authentication, configuration of SSL Certificates, and other data that should not be stored in plaintext configuration. Gloo runs an independent (goroutine) controller to monitor secrets. Secrets are stored in their own secret storage layer.'
                }
                exploreMoreLink={{
                  prompt: 'View Secrets',
                  link: '/settings/secrets/'
                }}
                descriptionMinHeight={'65px'}>
                <>
                  {/* <TallyInformationDisplay
                    tallyCount={0}
                    tallyDescription={`proxy error`}
                    color='orange'
                    moreInfoLink={{
                      prompt: 'View',
                      link: '/admin/proxy/?status=Rejected'
                    }}
                  /> */}
                  <GoodStateCongratulations typeOfItem={'secrets'} />

                  <TallyInformationDisplay
                    tallyCount={secretsList.length}
                    tallyDescription={`secret${
                      secretsList.length > 1 ? 's' : ''
                    }`}
                    color='blue'
                  />
                </>
              </StatusTile>
            </div>
            <div></div>
          </div>
        </div>
      </Container>
    </ErrorBoundary>
  );
};
