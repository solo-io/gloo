import * as React from 'react';
import { SoloTable } from 'Components/Common/SoloTable';

const routeColumns = [
  {
    title: 'Matcher',
    dataIndex: 'matcher',
    width: 200
  },
  {
    title: 'Path Match Type',
    dataIndex: 'pathMatch'
  },
  {
    title: 'Methods',
    dataIndex: 'method'
  },
  {
    title: 'Upstream',
    dataIndex: 'upstreamName'
  },
  {
    title: 'Destination',
    dataIndex: 'destinationName'
  },
  {
    title: 'Headers',
    dataIndex: 'header'
  },
  {
    title: 'Query Parameters',
    dataIndex: 'queryParams'
  },
  {
    title: 'Actions',
    dataIndex: 'actions',
    render: (text: any) => <div>ACTION!</div>
  }
];

const routeData: any[] = [];
for (let i = 1; i <= 5; i++) {
  routeData.push({
    key: i,
    matcher: '/test',
    pathMatch: 'PREFIX',
    method: '*',
    upstreamName: 'fake-upstream-13-9080',
    destinationName: 'test',
    header: 'test',
    queryParams: 'test',
    actions: ''
  });
}

export const Routes = () => {
  return (
    <div>
      <SoloTable columns={routeColumns} dataSource={routeData} />
    </div>
  );
};
