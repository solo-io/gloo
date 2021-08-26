import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { GlooAdminGateways } from './GlooAdminGateways';
import { GlooAdminProxy } from './GlooAdminProxy';
import { GlooAdminSettings } from './GlooAdminSettings';
import useSWR from 'swr';
import { glooInstanceApi } from 'API/gloo-instance';
import { GlooInstance } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/glooinstance_pb';
import { useListGlooInstances } from 'API/hooks';
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
  const { adminPage, name, namespace } = useParams();
  const navigate = useNavigate();

  const { data: glooInstances, error: instanceError } = useListGlooInstances();

  const [glooInstance, setGlooInstance] = useState<GlooInstance.AsObject>();

  useEffect(() => {
    if (!!glooInstances) {
      setGlooInstance(
        glooInstances.find(
          instance =>
            instance.metadata?.name === name &&
            instance.metadata?.namespace === namespace
        )
      );
    } else {
      setGlooInstance(undefined);
    }
  }, [name, namespace, glooInstances]);

  if (!!instanceError) {
    return <DataError error={instanceError} />;
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
