import styled from '@emotion/styled/macro';
import { graphqlApi } from 'API/graphql';
import { useIsGlooFedEnabled, useListGraphqlSchemas } from 'API/hooks';
import { ReactComponent as DownloadIcon } from 'assets/download-icon.svg';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { ReactComponent as GrpcIcon } from 'assets/grpc-icon.svg';
import { ReactComponent as RESTIcon } from 'assets/openapi-icon.svg';
import { ReactComponent as XIcon } from 'assets/x-icon.svg';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import { DataError } from 'Components/Common/DataError';
import ErrorModal from 'Components/Common/ErrorModal';
import { Loading } from 'Components/Common/Loading';
import { SectionCard } from 'Components/Common/SectionCard';
import { CheckboxFilterProps } from 'Components/Common/SoloCheckbox';
import { RenderSimpleLink, SimpleLinkProps } from 'Components/Common/SoloLink';
import {
  RenderStatus,
  SoloTable,
  TableActionCircle,
  TableActions,
} from 'Components/Common/SoloTable';
import { doDownload } from 'download-helper';
import { GraphqlSchema } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import React from 'react';
import { useNavigate, useParams } from 'react-router';
import { colors } from 'Styles/colors';
import { useDeleteAPI } from 'utils/hooks';
import { APIType } from './GraphqlLanding';
import { NewApiModal } from './NewApiModal';

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
  cluster: string;
  status: number;
  resolvers: number;
  actions: GraphqlSchema.AsObject;
};

export const GraphqlTable = (props: Props & TableHolderProps) => {
  const { name, namespace } = useParams();
  const { data: glooFedCheckResponse, error: glooFedCheckError } =
    useIsGlooFedEnabled();
  const isGlooFedEnabled = glooFedCheckResponse?.enabled;

  const {
    data: graphqlSchemas,
    error: graphqlSchemaError,
    mutate,
  } = useListGraphqlSchemas();

  const {
    isDeleting,
    triggerDelete,
    cancelDelete,
    closeErrorModal,
    errorModalIsOpen,
    errorDeleteModalProps,
    deleteFn,
  } = useDeleteAPI({
    revalidate: mutate,
    optimistic: true,
  });

  const [tableData, setTableData] = React.useState<TableDataType[]>([]);

  React.useEffect(() => {
    if (graphqlSchemas) {
      setTableData(
        graphqlSchemas.map(gqlSchema => {
          return {
            key: gqlSchema.metadata?.uid!,
            name: {
              displayElement: gqlSchema.metadata?.name ?? '',
              link: gqlSchema.metadata
                ? isGlooFedEnabled
                  ? `/gloo-instances/${gqlSchema.glooInstance?.namespace}/${gqlSchema.glooInstance?.name}/apis/${gqlSchema.metadata.clusterName}/${gqlSchema.metadata.namespace}/${gqlSchema.metadata.name}/`
                  : `/gloo-instances/${gqlSchema.glooInstance?.namespace}/${gqlSchema.glooInstance?.name}/apis/${gqlSchema.metadata.namespace}/${gqlSchema.metadata.name}/`
                : '',
            },
            namespace: gqlSchema.metadata?.namespace ?? '',
            cluster: gqlSchema.metadata?.clusterName ?? '',
            status: gqlSchema.status?.state ?? 0,
            resolvers:
              gqlSchema?.spec?.executableSchema?.executor?.local?.resolutionsMap
                ?.length ?? 0,
            actions: {
              ...gqlSchema,
            },
          };
        })
      );
    } else {
      setTableData([]);
    }
  }, [!!graphqlSchemas, graphqlSchemas?.length, isGlooFedEnabled]);

  const onDownloadSchema = (gqlSchema: GraphqlSchema.AsObject) => {
    if (gqlSchema.metadata) {
      graphqlApi
        .getGraphqlSchemaYaml({
          name: gqlSchema.metadata.name,
          namespace: gqlSchema.metadata.namespace,
          clusterName: gqlSchema.metadata.clusterName,
        })
        .then(gqlSchemaYaml => {
          doDownload(
            gqlSchemaYaml,
            gqlSchema.metadata?.namespace +
              '--' +
              gqlSchema.metadata?.name +
              '.yaml'
          );
        });
    }
  };
  if (!!graphqlSchemaError) {
    return <DataError error={graphqlSchemaError} />;
  } else if (!graphqlSchemas) {
    return <Loading message={'Retrieving GraphQL schemas...'} />;
  }

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
      dataIndex: 'resolvers',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      render: RenderStatus,
    },

    {
      title: 'Actions',
      dataIndex: 'actions',
      render: (gqlSchema: GraphqlSchema.AsObject) => (
        <TableActions className='space-x-3 '>
          <TableActionCircle onClick={() => onDownloadSchema(gqlSchema)}>
            <DownloadIcon />
          </TableActionCircle>
          <TableActionCircle
            onClick={() =>
              triggerDelete({
                name: gqlSchema.metadata?.name!,
                namespace: gqlSchema.metadata?.namespace!,
                clusterName: gqlSchema.metadata?.clusterName!,
              })
            }>
            <XIcon />
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
        flatTopped
        removeShadows
      />
      <ConfirmationModal
        visible={isDeleting}
        confirmPrompt='delete this API'
        confirmButtonText='Delete'
        goForIt={deleteFn}
        cancel={cancelDelete}
        isNegative
      />
      <ErrorModal
        {...errorDeleteModalProps}
        cancel={closeErrorModal}
        visible={errorModalIsOpen}
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
  // TODO:  Temporarily hides the Grpc and Rest boxes.  Remove when they are available.
  function hideBox(filter: CheckboxFilterProps) {
    switch (filter.label) {
      case APIType.GRAPHQL:
        return true;
      case APIType.REST:
        return false;
      case APIType.GRPC:
        return false;
      default:
        return false;
    }
  }
  const toggleGraphqlModal = () => {
    setShowGraphqlModal(!showGraphqlModal);
  };

  return (
    <>
      {typeFilters
        ?.filter(filter => {
          if (typeFilters?.some(f => f.checked)) {
            return filter.checked;
          }
          return hideBox(filter);
        })
        ?.map(filter => (
          <SectionCard
            key={filter.label}
            cardName={filter.label}
            logoIcon={<GraphqlIconHolder>{getIcon(filter)}</GraphqlIconHolder>}
            noPadding={true}>
            <GraphqlTable {...props} wholePage={true} />
          </SectionCard>
        ))}
      <NewApiModal
        showNewModal={showGraphqlModal}
        toggleNewModal={toggleGraphqlModal}
      />
    </>
  );
};
