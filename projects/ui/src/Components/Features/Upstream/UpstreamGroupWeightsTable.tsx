import React, { useEffect } from 'react';
import { useParams } from 'react-router';
import styled from '@emotion/styled/macro';
import { SoloTable } from 'Components/Common/SoloTable';
import { Loading } from 'Components/Common/Loading';
import { WeightPercentageBlock } from 'Styles/StyledComponents/weight-block';
import { colors } from 'Styles/colors';
import { ReactComponent as KubeIcon } from 'assets/kubernetes-icon.svg';
import { WeightedDestination } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/proxy_pb';
import {
  getWeightedDestinationType,
  TYPE_CONSUL,
  TYPE_KUBE,
  TYPE_STATIC,
  TYPE_OTHER,
} from 'utils/upstream-helpers';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { RenderSimpleLink } from 'Components/Common/SoloLink';
import { useIsGlooFedEnabled } from 'API/hooks';

/*
!!** NOTE: all of this Type stuff will likely come back, but is removed for immediate future.
  type DestinationTypeOptions =
  | typeof TYPE_CONSUL
  | typeof TYPE_KUBE
  | typeof TYPE_STATIC
  | typeof TYPE_OTHER;*/

const TableHolder = styled.div`
  table thead.ant-table-thead tr th {
    background: ${colors.marchGrey};
  }
`;

/*const TypeHolder = styled.div`
  display: flex;

  svg {
    margin-right: 8px;
  }
`;*/

type WeightTableFields = {
  key: string;
  weight: number;
  name: string;
  namespace: string;
  //type: DestinationTypeOptions;
};

interface Props {
  destinations: WeightedDestination.AsObject[] | undefined;
}

export const UpstreamGroupWeightsTable = ({ destinations }: Props) => {
  const [tableData, setTableData] = React.useState<WeightTableFields[]>([]);

  const {
    name: glooInstanceName,
    namespace: glooInstanceNamespace,
    upstreamGroupClusterName: cluster,
  } = useParams();

  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  useEffect(() => {
    if (destinations) {
      setTableData(
        destinations.map(wd => {
          return {
            key:
              wd.destination?.upstream?.name ??
              wd.destination?.kube?.ref?.name ??
              wd.destination?.consul?.serviceName ??
              'A destination was provided with no name',
            weight: wd.weight,
            name:
              wd.destination?.upstream?.name ??
              wd.destination?.kube?.ref?.name ??
              wd.destination?.consul?.serviceName ??
              '',
            namespace:
              wd.destination?.upstream?.namespace ??
              wd.destination?.kube?.ref?.namespace ??
              '',
            // type: getWeightedDestinationType(wd),
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [destinations]);

  if (!destinations) {
    return <Loading message={'Retrieving weights...'} />;
  }

  const renderWeight = (weight: number) => {
    const percentage = Math.round(
      (100 * weight) / destinations.reduce((acc, wd) => acc + wd.weight, 0)
    );
    return (
      <WeightPercentageBlock percentage={percentage} width='75px'>
        {percentage}%
      </WeightPercentageBlock>
    );
  };

  const renderName = (rowData?: any) => {
    const name = rowData.name;
    const namespace = rowData.namespace;
    if (
      name &&
      namespace &&
      cluster &&
      glooInstanceName &&
      glooInstanceNamespace
    ) {
      return (
        <RenderSimpleLink
          displayElement={name}
          link={
            isGlooFedEnabled
              ? `/gloo-instances/${glooInstanceNamespace}/${glooInstanceName}/upstreams/${cluster}/${namespace}/${name}`
              : `/gloo-instances/${glooInstanceNamespace}/${glooInstanceName}/upstreams/${namespace}/${name}`
          }
          inline
        />
      );
    }
    return name;
  };

  /*const renderType = (type: DestinationTypeOptions) => {
    return (
      <TypeHolder>
        <IconHolder width={25}>
          {
            type === 'Consul' ? (
            <ConsulIcon />
          ) :  type ===
            'Kubernetes' ? (
              <KubeIcon />  : type === 'Static' ? (
            <StaticIcon />
          )
            ) : null
          }
        </IconHolder>
        {type}
      </TypeHolder>
    );
  };*/

  let columns: any = [
    {
      title: 'Weight',
      dataIndex: 'weight',
      render: renderWeight,
    },
    {
      title: 'Name',
      render: renderName,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    /*{
      title: 'Type',
      dataIndex: 'type',
      render: renderType,
    },*/
  ];

  return (
    <TableHolder>
      <SoloTable
        columns={columns}
        dataSource={tableData}
        removePaging
        removeShadows
        curved={true}
      />
    </TableHolder>
  );
};
