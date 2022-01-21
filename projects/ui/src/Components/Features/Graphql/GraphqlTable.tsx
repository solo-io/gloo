import styled from '@emotion/styled/macro';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as FailoverIcon } from 'assets/GlooFed-Specific/failover-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as GrpcIcon } from 'assets/grpc-icon.svg';
import { ReactComponent as RESTIcon } from 'assets/openapi-icon.svg';
import { SectionCard } from 'Components/Common/SectionCard';
import { CheckboxFilterProps } from 'Components/Common/SoloCheckbox';
import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import { RenderSimpleLink, SimpleLinkProps } from 'Components/Common/SoloLink';
import {
  RenderCluster,
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import {
  ExecutableSchema,
  GraphQLRouteConfig,
} from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/external/envoy/extensions/graphql/graphql_pb';
import { Upstream } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/gloo_resources_pb';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { APIType } from './GraphqlLanding';
import bookInfoSchema from './data/book-info.json';
import petstoreSchema from './data/petstore.json';
import { NewApiModal } from './NewApiModal';
import { ApiProvider } from './state/ApiProvider.state';

export const GraphqlIconHolder = styled.div`
  display: flex;
  align-items: center;
  justify-items: center;

  svg {
    width: 35px;
    max-width: none;
  }
`;

const PositionHolder = styled.div`
  position: relative;
`;

const SecondaryComponent = styled.div`
  position: absolute;
  right: 20px;
`;

const Button = styled.button`
  color: ${colors.oceanBlue};
  &:hover {
    cursor: pointer;
    color: ${colors.seaBlue};
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
  actions: typeof bookInfoSchema.spec;
};

let testData: TableDataType[] = [
  {
    key: bookInfoSchema.metadata.uid,
    name: {
      displayElement: bookInfoSchema.metadata.name,
      link: `/apis/${bookInfoSchema.metadata.namespace}/${bookInfoSchema.metadata.name}`,
    },
    namespace: bookInfoSchema.metadata.namespace,
    glooInstance: {
      name: 'local-gloo-system',
      namespace: 'gloo-system',
    },
    cluster: 'local',
    failover: false,
    status: 1,
    actions: {
      ...bookInfoSchema.spec,
    },
  },
];

export const GraphqlTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const navigate = useNavigate();
  const [tableData, setTableData] = React.useState<TableDataType[]>(testData);

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
      title: 'Resolvers',
      dataIndex: 'failover',
      render: () => <div>6</div>,
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
  const [showGraphqlModal, setShowGraphqlModal] = React.useState(false);
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
  const toggleGraphqlModal = () => {
    setShowGraphqlModal(!showGraphqlModal);
  }

  function getLink(filter: CheckboxFilterProps) {
    return filter.label === APIType.GRAPHQL ? (
    <SecondaryComponent>
      <Button type="button" onClick={toggleGraphqlModal} className=''>
         Create Graphql API
      </Button>

    </SecondaryComponent>) : null;
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
            secondaryComponent={getLink(filter)}
            >
              <GraphqlTable {...props} wholePage={true} />
            </SectionCard>
        ))}
        <ApiProvider>
          <NewApiModal showNewModal={showGraphqlModal} toggleNewModal={toggleGraphqlModal} />
        </ApiProvider>
    </>
  );
};
