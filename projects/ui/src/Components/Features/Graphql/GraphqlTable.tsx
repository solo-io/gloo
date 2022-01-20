import styled from '@emotion/styled/macro';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as FailoverIcon } from 'assets/GlooFed-Specific/failover-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as GrpcIcon } from 'assets/grpc-icon.svg';
import { ReactComponent as RESTIcon } from 'assets/openapi-icon.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { CheckboxFilterProps } from 'Components/Common/SoloCheckbox';
import { RenderSimpleLink, SimpleLinkProps } from 'Components/Common/SoloLink';
import {
  RenderCluster,
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { APIType } from './GraphqlLanding';

const GraphqlIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

type TableHolderProps = { wholePage?: boolean };
const TableHolder = styled.div<TableHolderProps>`
  ${(props: TableHolderProps) =>
    props.wholePage
      ? ''
      : `
    table thead.ant-table-thead tr th {
      background: ${colors.marchGrey};
    }
  `};
`;

type Props = {
  typeFilters?: CheckboxFilterProps[];
};

type TableDataType = {
  key: string;
  name: SimpleLinkProps;
  namespace: string;
  glooInstance?: { name: string; namespace: string };
  cluster: string;
  failover: boolean;
  status: number;
  // TODO: replace this with mock graphql data
  actions: Upstream.AsObject;
};
let testData: TableDataType[] = [
  {
    key: 'd65e0c8a-5e7b-4788-adef-259b677d7982',
    name: {
      displayElement: 'default-details-9080',
      link: '/gloo-instances/gloo-system/local-gloo-system/upstreams/local/gloo-system/default-details-9080',
    },
    namespace: 'gloo-system',
    glooInstance: {
      name: 'local-gloo-system',
      namespace: 'gloo-system',
    },
    cluster: 'local',
    failover: false,
    status: 1,
    actions: {
      metadata: {
        name: 'default-details-9080',
        namespace: 'gloo-system',
        uid: 'd65e0c8a-5e7b-4788-adef-259b677d7982',
        resourceVersion: '5731698',
        labelsMap: [['discovered_by', 'kubernetesplugin']],
        annotationsMap: [
          ['cloud.google.com/neg', '{"ingress":true}'],
          [
            'kubectl.kubernetes.io/last-applied-configuration',
            '{"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"app":"details","service":"details"},"name":"details","namespace":"default"},"spec":{"ports":[{"name":"http","port":9080}],"selector":{"app":"details"}}}\n',
          ],
        ],
        clusterName: 'local',
      },
      spec: {
        discoveryMetadata: {
          labelsMap: [
            ['app', 'details'],
            ['service', 'details'],
          ],
        },
        healthChecksList: [],
        kube: {
          serviceName: 'details',
          serviceNamespace: 'default',
          servicePort: 9080,
          selectorMap: [['app', 'details']],
        },
      },
      status: {
        state: 1,
        reason: '',
        reportedBy: 'gloo',
        subresourceStatusesMap: [],
      },
      glooInstance: {
        name: 'local-gloo-system',
        namespace: 'gloo-system',
      },
    },
  },
  {
    key: 'd20e5111-bc6c-4363-b1d9-0f92008ed8ed',
    name: {
      displayElement: 'default-kubernetes-443',
      link: '/gloo-instances/gloo-system/local-gloo-system/upstreams/local/gloo-system/default-kubernetes-443',
    },
    namespace: 'gloo-system',
    glooInstance: {
      name: 'local-gloo-system',
      namespace: 'gloo-system',
    },
    cluster: 'local',
    failover: false,
    status: 1,
    actions: {
      metadata: {
        name: 'default-kubernetes-443',
        namespace: 'gloo-system',
        uid: 'd20e5111-bc6c-4363-b1d9-0f92008ed8ed',
        resourceVersion: '2559',
        labelsMap: [['discovered_by', 'kubernetesplugin']],
        annotationsMap: [],
        clusterName: 'local',
      },
      spec: {
        discoveryMetadata: {
          labelsMap: [
            ['component', 'apiserver'],
            ['provider', 'kubernetes'],
          ],
        },
        healthChecksList: [],
        kube: {
          serviceName: 'kubernetes',
          serviceNamespace: 'default',
          servicePort: 443,
          selectorMap: [],
        },
      },
      status: {
        state: 1,
        reason: '',
        reportedBy: 'gloo',
        subresourceStatusesMap: [],
      },
      glooInstance: {
        name: 'local-gloo-system',
        namespace: 'gloo-system',
      },
    },
  },
  {
    key: '479e99b2-3ba6-4197-a45b-ed2f075a9fac',
    name: {
      displayElement: 'default-kubernetes-443',
      link: '/gloo-instances/gloo-system/remote-gloo-system/upstreams/remote/gloo-system/default-kubernetes-443',
    },
    namespace: 'gloo-system',
    glooInstance: {
      name: 'remote-gloo-system',
      namespace: 'gloo-system',
    },
    cluster: 'remote',
    failover: false,
    status: 1,
    actions: {
      metadata: {
        name: 'default-kubernetes-443',
        namespace: 'gloo-system',
        uid: '479e99b2-3ba6-4197-a45b-ed2f075a9fac',
        resourceVersion: '3403',
        labelsMap: [['discovered_by', 'kubernetesplugin']],
        annotationsMap: [],
        clusterName: 'remote',
      },
      spec: {
        discoveryMetadata: {
          labelsMap: [
            ['component', 'apiserver'],
            ['provider', 'kubernetes'],
          ],
        },
        healthChecksList: [],
        kube: {
          serviceName: 'kubernetes',
          serviceNamespace: 'default',
          servicePort: 443,
          selectorMap: [],
        },
      },
      status: {
        state: 1,
        reason: '',
        reportedBy: 'gloo',
        subresourceStatusesMap: [],
      },
      glooInstance: {
        name: 'remote-gloo-system',
        namespace: 'gloo-system',
      },
    },
  },
];

export const GraphqlTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();
  const [tableData, setTableData] = React.useState<TableDataType[]>(testData);
  const multipleClustersOrInstances = false;

  const renderGlooInstanceList = (glooInstance: {
    name: string;
    namespace: string;
  }) => (
    <div
      onClick={() =>
        navigate(
          `/gloo-instances/${glooInstance.namespace}/${glooInstance.name}/`
        )
      }
    >
      {glooInstance.name}
    </div>
  );

  const renderFailover = (failoverExists: boolean) => {
    return failoverExists ? (
      <IconHolder>
        <FailoverIcon />
      </IconHolder>
    ) : (
      <React.Fragment />
    );
  };

  let columns: any = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: RenderSimpleLink,
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Failover',
      dataIndex: 'failover',
      render: renderFailover,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (upstream: Upstream.AsObject) => (
        <TableActions>
          <TableActionCircle onClick={() => {}}>
            <DownloadIcon />
          </TableActionCircle>
        </TableActions>
      ),
    },
  ];

  if (props.wholePage && multipleClustersOrInstances) {
    columns.splice(2, 0, {
      title: 'Cluster',
      dataIndex: 'cluster',
      render: RenderCluster,
    });
  }
  if (props.wholePage) {
    columns.splice(2, 0, {
      title: 'Gloo Instance',
      dataIndex: 'glooInstance',
      render: renderGlooInstanceList,
    });
  }

  return (
    <TableHolder wholePage={props.wholePage}>
      <SoloTable
        columns={columns}
        dataSource={tableData}
        removePaging
        removeShadows
        curved={props.wholePage}
      />
    </TableHolder>
  );
};

export const GraphqlPageTable = (props: Props) => {
  const { typeFilters } = props;
  function getIcon(filter: CheckboxFilterProps) {
    switch (filter.label) {
      case APIType.GRAPHQL:
        return <GraphQLIcon />;
      case APIType.REST:
        return <RESTIcon />;
      case APIType.GRPC:
        return <GrpcIcon />;
      default:
        break;
    }
  }
  return (
    <>
      {typeFilters
        ?.filter(filter => {
          if (typeFilters?.some(f => f.checked)) {
            return filter.checked;
          }
          return true;
        })
        ?.map(filter => (
          <SectionCard
            key={filter.label}
            cardName={filter.label}
            logoIcon={<GraphqlIconHolder>{getIcon(filter)}</GraphqlIconHolder>}
            noPadding={true}
          >
            <GraphqlTable {...props} wholePage={true} />
          </SectionCard>
        ))}
    </>
  );
};
