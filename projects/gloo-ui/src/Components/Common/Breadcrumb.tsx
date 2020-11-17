import styled from '@emotion/styled';
import { Breadcrumb as AntdBreadcrumb } from 'antd';
import * as React from 'react';
import { useLocation } from 'react-router';
import { Link } from 'react-router-dom';

const BreadcrumbContainer = styled(AntdBreadcrumb)`
  margin-bottom: 15px;
`;

const breadcrumbNameMap: { [key: string]: string } = {
  '/overview': 'Overview',
  '/virtualservices': 'Virtual Services',
  '/routetables': 'Route Tables',
  '/upstreams': 'Upstreams',
  '/upstreams/upstreamgroups': 'Upstream Groups',
  '/admin': 'Admin',
  '/settings': 'Settings',
  '/settings/secrets': 'Secrets',
  '/wasm': 'Wasm'
};
export interface RouteParams {
  virtualservicename?: string;
  routetablename?: string;
  settingsublocation?: string; // This is currently unused as Settings doesn't  have the same flow.
  sublocation?: string;
}

export const Breadcrumb = () => {
  let location = useLocation();

  // https://ant.design/components/breadcrumb/
  let pathSnippets = location.pathname.split('/').filter(i => i);

  let extraBreadcrumbItems = pathSnippets.map((_, index) => {
    const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;

    if (breadcrumbNameMap[url] !== undefined) {
      return (
        <AntdBreadcrumb.Item key={url}>
          <Link to={url}>{breadcrumbNameMap[url]}</Link>
        </AntdBreadcrumb.Item>
      );
    } else {
      return (
        <AntdBreadcrumb.Item key={url}>
          <Link to={url}>{pathSnippets[pathSnippets.length - 1]}</Link>
        </AntdBreadcrumb.Item>
      );
    }
  });

  // do not make a breadcrumb for the namespace part of the url
  if (pathSnippets.length >= 3) {
    delete extraBreadcrumbItems[pathSnippets.length - 2];
  }

  const breadcrumbItems = [
    <AntdBreadcrumb.Item key='home'>
      <Link to='/'>Home</Link>
    </AntdBreadcrumb.Item>
  ].concat(extraBreadcrumbItems);

  return (
    <BreadcrumbContainer separator='>'>
      {extraBreadcrumbItems}
    </BreadcrumbContainer>
  );
};
