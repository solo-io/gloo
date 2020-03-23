import styled from '@emotion/styled';
import { Breadcrumb } from 'Components/Common/Breadcrumb';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { SectionCard } from 'Components/Common/SectionCard';
import React from 'react';
import { configAPI } from 'store/config/api';
import useSWR, { mutate } from 'swr';
import { ErrorBoundary } from '../Errors/ErrorBoundary';

const Heading = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`;
export const Settings = () => {
  const { data: settingsDetails, error: settingsError } = useSWR(
    'getSettings',
    configAPI.getSettings
  );

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
  if (!settingsDetails) return <div>Loading...</div>;
  return (
    <ErrorBoundary
      fallback={<div>There was an error with the Settings section</div>}>
      <div>
        <SectionCard cardName={'Settings YAML'}>
          <ConfigDisplayer
            content={settingsDetails?.raw?.content || ''}
            asEditor
            yamlError={settingsError}
            saveEdits={handleSaveYAML}
          />
        </SectionCard>
      </div>
    </ErrorBoundary>
  );
};
