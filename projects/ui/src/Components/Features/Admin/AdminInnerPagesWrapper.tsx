import React, { useEffect } from 'react';
import { useParams, useNavigate, Routes, Route } from 'react-router';
import { colors } from 'Styles/colors';
import styled from '@emotion/styled';
import { SoloRadioGroup } from 'Components/Common/SoloRadioGroup';
import { FederatedVirtualServices } from './FederatedVirtualServices';
import { FederatedUpstreams } from './FederatedUpstreams';
import { FederatedUpstreamGroups } from './FederatedUpstreamGroups';
import { FederatedSettingsTable } from './FederatedSettingsTable';
import { FederatedRouteTables } from './FederatedRouteTables';
import { FederatedGateways } from './FederatedGateways';
import { FederatedAuthorizedConfigurations } from './FederatedAuthorizedConfigurations';
import { FederatedRateLimits } from './FederatedRateLimits';

const AdminInnerPagesContainer = styled.div`
  display: grid;
  grid-template-columns: 200px 1fr;
  grid-gap: 28px;
`;

const pageOptions: {
  displayName: string;
  id: string | number;
}[] = [
  {
    displayName: 'Virtual Services',
    id: 'virtual-services',
  },
  {
    displayName: 'Upstreams',
    id: 'upstreams',
  },
  {
    displayName: 'Upstream Groups',
    id: 'upstream-groups',
  },
  {
    displayName: 'Auth Configs',
    id: 'authorizations',
  },
  {
    displayName: 'Rate Limits',
    id: 'rate-limits',
  },
  {
    displayName: 'Route Tables',
    id: 'route-tables',
  },
  {
    displayName: 'Gateways',
    id: 'gateways',
  },
  {
    displayName: 'Settings',
    id: 'settings',
  },
];

export const AdminInnerPagesWrapper = () => {
  const { adminPage } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    if (!adminPage?.length) {
      navigate('virtual-services/');
    }
    /* eslint-disable-next-line react-hooks/exhaustive-deps */
  }, [adminPage]);

  const goToPage = (newPage: string | number | undefined) => {
    navigate(`../${newPage as string}/`);
  };

  return (
    <AdminInnerPagesContainer>
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
          adminPage === 'virtual-services' ? (
            <FederatedVirtualServices />
          ) : adminPage === 'upstreams' ? (
            <FederatedUpstreams />
          ) : adminPage === 'upstream-groups' ? (
            <FederatedUpstreamGroups />
          ) : adminPage === 'authorizations' ? (
            <FederatedAuthorizedConfigurations />
          ) : adminPage === 'rate-limits' ? (
            <FederatedRateLimits />
          ) : adminPage === 'route-tables' ? (
            <FederatedRouteTables />
          ) : adminPage === 'gateways' ? (
            <FederatedGateways />
          ) : adminPage === 'settings' ? (
            <FederatedSettingsTable />
          ) : (
            'This page does not exist.'
          )
        }
      </div>
    </AdminInnerPagesContainer>
  );
};
