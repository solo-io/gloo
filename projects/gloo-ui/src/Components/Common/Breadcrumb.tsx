import styled from '@emotion/styled';
import { Breadcrumb as AntdBreadcrumb } from 'antd';
import * as React from 'react';
import { useHistory, useLocation, useParams } from 'react-router';
import { colors } from 'Styles';

const BreadcrumbContainer = styled(AntdBreadcrumb)`
  margin-bottom: 15px;
`;

type CrumbLinkProps = { clickable?: boolean };
const CrumbLink = styled.span<CrumbLinkProps>`
  font-size: 14px;

  ${(props: CrumbLinkProps) =>
    props.clickable
      ? `
      color: ${colors.seaBlue};
      cursor: pointer;`
      : `
      color: ${colors.septemberGrey};`};
`;

const CapitalizedCrumbLink = styled(CrumbLink)`
  text-transform: capitalize;
`;

const rootNameMap: { [key: string]: string } = {
  virtualservices: 'Virtual Services',
  upstreams: 'Upstreams',
  stats: 'Stats',
  settings: 'Settings',
  admin: 'Administration'
};

export interface RouteParams {
  virtualservicename?: string;
  settingsublocation?: string; // This is currently unused as Settings doesn't  have the same flow.
  sublocation?: string;
}

export const Breadcrumb = () => {
  let location = useLocation();
  let history = useHistory();
  let { virtualservicename, sublocation } = useParams();

  const goToRoot = () => {
    history.push({
      pathname: `/${location.pathname.split('/')[1]}/`
    });
  };

  return (
    <BreadcrumbContainer separator='>'>
      <AntdBreadcrumb.Item onClick={goToRoot}>
        <CrumbLink clickable={true}>
          {rootNameMap[location.pathname.split('/')[1]] as string}
        </CrumbLink>
      </AntdBreadcrumb.Item>
      {!!location.search && (
        <AntdBreadcrumb.Item>
          <CrumbLink>{location.search.split('=')[1]}</CrumbLink>
        </AntdBreadcrumb.Item>
      )}
      {!!virtualservicename && (
        <AntdBreadcrumb.Item>
          <CrumbLink>{virtualservicename}</CrumbLink>
        </AntdBreadcrumb.Item>
      )}
      {!!sublocation && (
        <AntdBreadcrumb.Item>
          <CapitalizedCrumbLink>{sublocation}</CapitalizedCrumbLink>
        </AntdBreadcrumb.Item>
      )}
    </BreadcrumbContainer>
  );
};
