import { LinkOutlined } from '@ant-design/icons';
import { Global } from '@emotion/core';
import { ColumnsType } from 'antd/lib/table';
import {
  useGetGraphqlApiDetails,
  useGetGraphqlApiYaml,
  useGetStitchedSchemaDefinition,
  useIsGlooFedEnabled,
  usePageGlooInstance,
} from 'API/hooks';
import { SoloInput } from 'Components/Common/SoloInput';
import { SoloTable } from 'Components/Common/SoloTable';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1beta1/graphql_pb';
import React, { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { makeGraphqlApiLink } from 'utils/graphql-helpers';
import StitchedGqlAddSubGraph from './update-sub-graph/StitchedGqlAddSubGraph';
import StitchedGqlRemoveSubGraph from './update-sub-graph/StitchedGqlRemoveSubGraph';
import { styles } from './StitchedGqlSubGraphs.style';
import StitchedGqlEditTypeMergeMap from './update-sub-graph/StitchedGqlEditTypeMergeMap';

const StitchedGqlSubGraphs: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
  subGraphs: StitchedSchema.SubschemaConfig.AsObject[];
}> = ({ apiRef, subGraphs }) => {
  const { glooInstance } = usePageGlooInstance();
  const isGlooFedEnabled = useIsGlooFedEnabled().data?.enabled;
  const { data: graphqlApi, mutate: mutateDetails } =
    useGetGraphqlApiDetails(apiRef);
  const { mutate: mutateSchema } = useGetStitchedSchemaDefinition(apiRef);
  const { mutate: mutateYaml } = useGetGraphqlApiYaml(apiRef);
  const mutateAll = () => {
    mutateDetails();
    mutateSchema();
    mutateYaml();
  };

  // -- SEARCHING -- //
  const [searchText, setSearchText] = useState('');
  const [filteredData, setFilteredData] = useState(subGraphs);
  useEffect(() => {
    const lstext = searchText.toLowerCase();
    setFilteredData(
      subGraphs.filter(d => d.name.toLowerCase().includes(lstext))
    );
  }, [searchText, subGraphs]);

  // -- TABLE COLUMNS -- //
  const columns = useMemo(() => {
    return [
      {
        title: 'Name',
        dataIndex: 'name',
        width: 300,
        render: (
          value: any,
          record: StitchedSchema.SubschemaConfig.AsObject,
          index: number
        ) => {
          return (
            <Link
              to={makeGraphqlApiLink(
                record.name,
                record.namespace,
                glooInstance?.spec?.cluster,
                graphqlApi?.glooInstance?.name,
                graphqlApi?.glooInstance?.namespace,
                isGlooFedEnabled
              )}>
              <LinkOutlined />
              &nbsp;&nbsp;{record.name}
            </Link>
          );
        },
      },
      {
        title: 'Namespace',
        dataIndex: 'namespace',
      },
      {
        title: 'Actions',
        dataIndex: 'actions',
        render: (
          value: any,
          subGraphConfig: StitchedSchema.SubschemaConfig.AsObject,
          index: number
        ) => {
          return (
            <div className='flex align-center'>
              <StitchedGqlRemoveSubGraph
                onAfterRemove={mutateAll}
                subGraphConfig={subGraphConfig}
              />
              <StitchedGqlEditTypeMergeMap
                onAfterEdit={mutateAll}
                subGraphConfig={subGraphConfig}
              />
            </div>
          );
        },
      },
    ] as ColumnsType<any>;
  }, [graphqlApi, isGlooFedEnabled, glooInstance]);

  if (!graphqlApi) return null;
  return (
    <>
      <Global styles={styles} />
      <div className='text-lg mb-5'>Sub Graphs</div>
      <div className='flex flex-wrap w-100 mb-5 gap-[10px]'>
        <div className='w-[400px]'>
          <SoloInput
            value={searchText}
            placeholder='Search by name...'
            onChange={s => setSearchText(s.target.value)}
          />
        </div>
        <div className='grow flex justify-end'>
          <StitchedGqlAddSubGraph onAfterAdd={mutateAll} />
        </div>
      </div>

      <SoloTable
        columns={columns}
        dataSource={filteredData}
        rowKey='name'
        removePaging
      />
    </>
  );
};

export default StitchedGqlSubGraphs;
