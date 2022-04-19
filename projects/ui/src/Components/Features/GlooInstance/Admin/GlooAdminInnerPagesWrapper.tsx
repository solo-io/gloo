import React from 'react';
import { useParams, useNavigate } from 'react-router';
import styled from '@emotion/styled';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { GlooAdminGateways } from './GlooAdminGateways';
import { GlooAdminProxy } from './GlooAdminProxy';
import { GlooAdminSettings } from './GlooAdminSettings';
import { usePageGlooInstance } from 'API/hooks';
import { DataError } from 'Components/Common/DataError';
import { Loading } from 'Components/Common/Loading';
import { GlooAdminEnvoy } from './GlooAdminEnvoy';
import { GlooAdminWatchedNamespaces } from './GlooAdminWatchNamespaces';
import { GlooAdminSecrets } from './GlooAdminSecrets';

const GlooAdminInnerPagesContainer = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;

  .YamlDisplayerContainer {
    max-width: 1000px;
  }
`;

const pageOptions: {
  displayName: string;
  id: string | number;
}[] = [
  {
    displayName: 'Gateways',
    id: 'gateways',
  },
  {
    displayName: 'Proxy',
    id: 'proxy',
  },
  {
    displayName: 'Envoy',
    id: 'envoy',
  },
  {
    displayName: 'Settings',
    id: 'settings',
  },
  {
    displayName: 'Watched Namespaces',
    id: 'watched-namespaces',
  },
  /*{
    displayName: 'Secrets',
    id: 'secrets',
  },*/
];

export const GlooAdminInnerPagesWrapper = () => {
  const { adminPage, name } = useParams();
  const navigate = useNavigate();

  const { glooInstance, instancesError } = usePageGlooInstance();

  if (!!instancesError) {
    return <DataError error={instancesError} />;
  } else if (!glooInstance?.spec) {
    return <Loading message={`Retrieving instance ${name}...`} />;
  }

  const goToPage = (newPage: string | number | undefined) => {
    navigate(`../${newPage as string}/`);
  };

  return (
    <GlooAdminInnerPagesContainer>
      <SoloRadioGroup
        options={pageOptions}
        currentSelection={adminPage}
        onChange={goToPage}
        withoutCheckboxes={true}
        forceAChoice={true}
      />
      <div>
        {
          /* We could use React Router here, but I don't think we get anything out 
            of it. We could do in <Content />, but then lose the slickness for Radio above. */
          adminPage === 'gateways' ? (
            <GlooAdminGateways />
          ) : adminPage === 'proxy' ? (
            <GlooAdminProxy />
          ) : adminPage === 'envoy' ? (
            <GlooAdminEnvoy />
          ) : adminPage === 'settings' ? (
            <GlooAdminSettings />
          ) : adminPage === 'watched-namespaces' ? (
            <GlooAdminWatchedNamespaces glooInstance={glooInstance} />
          ) : adminPage === 'secrets' ? (
            <GlooAdminSecrets />
          ) : (
            'This page does not exist.'
          )
        }
      </div>
    </GlooAdminInnerPagesContainer>
  );
};
