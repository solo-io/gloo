import { Global } from '@emotion/core';
import { Collapse } from 'antd';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import SoloAddButton from 'Components/Common/SoloAddButton';
import { SoloInput } from 'Components/Common/SoloInput';
import SoloNoMatches from 'Components/Common/SoloNoMatches';
import { SoloTable } from 'Components/Common/SoloTable';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useState } from 'react';
import { styles } from './GatewayGraphqlSubGraphs.style';

const GatewayGraphqlSubGraphs: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const [searchText, setSearchText] = useState('');

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
          <SoloAddButton onClick={() => alert('Open modal')}>
            Add Sub Graph
          </SoloAddButton>
        </div>
      </div>

      {searchText !== '' ? (
        <>
          <hr />
          <SoloNoMatches />
        </>
      ) : (
        <Collapse className='mb-5' defaultActiveKey={[0]}>
          <Collapse.Panel
            key={0}
            className='p-0'
            header={
              <div className='inline font-medium text-gray-900 whitespace-nowrap'>
                <GraphQLIcon className='w-4 h-4 fill-current inline' />
                &nbsp;&nbsp;Sub Graph 1
              </div>
            }>
            {/* Sub Graph 1 */}
            <SoloTable
              columns={[
                {
                  title: 'Object Name',
                  dataIndex: 'name',
                  width: 200,
                  // render: RenderSimpleLink,
                },
                {
                  title: 'Namespace',
                  dataIndex: 'namespace',
                },
              ]}
              dataSource={[{ name: 'Test Field', namespace: 'gloo-default' }]}
              removeShadows
              removePaging
            />
          </Collapse.Panel>
          <Collapse.Panel
            key={1}
            header={
              <div className='inline font-medium text-gray-900 whitespace-nowrap'>
                <GraphQLIcon className='w-4 h-4 fill-current inline' />
                &nbsp;&nbsp;Sub Graph 2
              </div>
            }>
            {/* Sub Graph 2 */}
            <SoloTable
              columns={[
                {
                  title: 'Object Name',
                  dataIndex: 'name',
                  width: 200,
                  // render: RenderSimpleLink,
                },
                {
                  title: 'Namespace',
                  dataIndex: 'namespace',
                },
              ]}
              dataSource={[
                { name: 'Test Field', namespace: 'gloo-default' },
                { name: 'Test Field 2', namespace: 'gloo-default' },
              ]}
              removeShadows
              removePaging
            />
          </Collapse.Panel>
        </Collapse>
      )}
    </>
  );
};

export default GatewayGraphqlSubGraphs;
