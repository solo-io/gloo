import React, { useEffect, useState } from 'react';
import styled from '@emotion/styled/macro';
import { Loading } from 'Components/Common/Loading';
import { useParams } from 'react-router';
import { HealthNotificationBox } from 'Components/Common/HealthNotificationBox';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as RouteTableIcon } from 'assets/route-icon.svg';
import { useListRouteTables } from 'API/hooks';
import { RouteTable } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gateway_resources_pb';
import { DataError } from 'Components/Common/DataError';
import { RoutesTable } from './RouteTable/RoutesTable';

const GlooIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 25px;
  }
`;

export const RouteTableDetails = () => {
  const { name, namespace, routeTableName, routeTableNamespace } = useParams();

  const { data: allRouteTables, error: rtError } = useListRouteTables({
    name: name!,
    namespace: namespace!,
  });
  const [routeTable, setRouteTable] = useState<RouteTable.AsObject>();

  useEffect(() => {
    if (!!allRouteTables) {
      setRouteTable(
        allRouteTables.find(
          rt =>
            rt.metadata?.name === routeTableName &&
            rt.metadata?.namespace === routeTableNamespace
        )
      );
    } else {
      setRouteTable(undefined);
    }
  }, [name, namespace, allRouteTables, routeTableName, routeTableNamespace]);

  if (!!rtError) {
    return <DataError error={rtError} />;
  } else if (!allRouteTables) {
    return <Loading message={'Retrieving route tables...'} />;
  }

  return (
    <SectionCard
      cardName={routeTableName!}
      logoIcon={
        <GlooIconHolder>
          <RouteTableIcon />
        </GlooIconHolder>
      }
      headerSecondaryInformation={[
        {
          title: 'Namespace',
          value: routeTableNamespace,
        },
      ]}
      health={{
        state: routeTable?.status?.state ?? 0,
        reason: routeTable?.status?.reason,
      }}>
      <>
        <HealthNotificationBox
          state={routeTable?.status?.state}
          reason={routeTable?.status?.reason}
        />
        <RoutesTable routes={routeTable?.spec?.routesList || []} />
      </>
    </SectionCard>
  );
};
