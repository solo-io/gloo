import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { Breadcrumb as AntdBreadcrumb } from 'antd';
import { colors } from 'Styles';

const BreadcrumbContainer = styled(AntdBreadcrumb)`
  margin-bottom: 15px;
`;

const CrumbLink = styled<'span', { clickable?: boolean }>('span')`
  font-size: 14px;

  ${props =>
    props.clickable
      ? `
      color: ${colors.seaBlue};
      cursor: pointer;`
      : `
      color: ${colors.septemberGrey};`};
`;

const rootNameMap: { [key: string]: string } = {
  virtualservices: 'Virtual Services',
  upstreams: 'Upstreams',
  stats: 'Stats',
  settings: 'Settings'
};

export interface RouteParams {
  virtualservicename?: string;
  settingsublocation?: string; // This is currently unused as Settings doesn't  have the same flow.
}

function BreadcrumbC({
  history,
  match,
  location
}: RouteComponentProps<RouteParams>) {
  const goToRoot = () => {
    history.push({
      pathname: `/${match.path.split('/')[1]}/`
    });
  };

  return (
    <BreadcrumbContainer separator='>'>
      <AntdBreadcrumb.Item onClick={goToRoot}>
        <CrumbLink clickable={true}>
          {rootNameMap[match.path.split('/')[1]] as string}
        </CrumbLink>
      </AntdBreadcrumb.Item>
      {!!location.search && (
        <AntdBreadcrumb.Item>
          <CrumbLink>{location.search.split('=')[1]}</CrumbLink>
        </AntdBreadcrumb.Item>
      )}
      {!!match.params.virtualservicename && (
        <AntdBreadcrumb.Item>
          <CrumbLink>{match.params.virtualservicename}</CrumbLink>
        </AntdBreadcrumb.Item>
      )}
    </BreadcrumbContainer>
  );
}

export const Breadcrumb = withRouter(BreadcrumbC);
