import React from 'react';
import styled from '@emotion/styled';
import { useParams, useLocation } from 'react-router';
import { Breadcrumb as AntdBreadcrumb } from 'antd';
import { colors } from 'Styles/colors';
import { Link } from 'react-router-dom';
import { AppName } from 'Components/Common/AppName';

type FixedFloaterProps = {
  isFloating: boolean;
};
const Container = styled.div<FixedFloaterProps>`
  display: flex;
  align-items: center;
  padding: 30px 0 0;
  width: 1275px;
  max-width: 100vw;
  margin: 0 auto;

  ${(props: FixedFloaterProps) =>
    props.isFloating && `${BreadcrumbContainer} {position: fixed;}`}

  .ant-breadcrumb > span:last-of-type {
    .ant-breadcrumb-link,
    .ant-breadcrumb-link a {
      cursor: default;
      color: ${colors.septemberGrey};
      pointer-events: none;
      font-weight: 400;
    }
    .ant-breadcrumb-separator {
      display: none;
    }
  }

  .ant-breadcrumb-separator {
    font-size: 14px;
    margin-right: 5px;
    color: ${colors.septemberGrey};
    pointer-events: none;
  }
`;

const BreadcrumbContainer = styled(AntdBreadcrumb)`
  cursor: default;
`;

const CrumbLink = styled.span`
  a {
    color: ${colors.seaBlue};
    font-size: 16px;
    cursor: pointer;
    text-decoration: none;
    margin-right: 5px;

    /*&:hover,
    &:active {
      text-decoration: underline;
    }*/
  }
`;

const breadcrumbLinkBackNames: { [key: string]: string | React.ReactNode } = {
  '/gloo-instances': 'Gloo Instances',
  '/virtual-services': 'Virtual Services',
  '/upstreams': 'Upstreams',
  '/upstream-groups': 'Upstream Groups',
  '/admin': (
    <>
      <AppName /> Administration
    </>
  ),
};

const breadcrumbNameMap: { [key: string]: string } = {
  'virtual-services': 'Virtual Services',
  upstreams: 'Upstreams',
  'upstream-groups': 'Upstream Groups',
  authorizations: 'Authorization',
  'route-tables': 'Route Tables',
  gateways: 'Gateways',
  settings: 'Settings',
  proxy: 'Proxy',
  'watched-namespaces': 'Watched Namespaces',
  'wasm-filters': 'Wasm Filters',
  apis: 'APIs',
};

const skippableNames: { [key: string]: boolean } = {
  '/gloo-instances/details': true,
};

const breakPathNames: { [key: string]: string } = {
  'virtual-services': 'Virtual Services',
  'route-tables': 'Route Tables',
  upstreams: 'Upstreams',
  'upstream-groups': 'Upstream Groups',
};

const nonLinkNames: { [key: string]: string } = {
  'federated-resources': 'Federated Resources',
};

const pathsHaveClutterAfter: { [key: string]: number[] } = {
  'gloo-instances': [1], // namespace
  'virtual-services': [1, 2], // clustername + namespace
  'route-tables': [1, 2], // clustername + namespace
  upstreams: [1, 2], // clustername + namespace
  'upstream-groups': [1, 2], // clustername + namespace
};

export interface RouteParams {
  name?: string;
  namespace?: string;
}

export function Breadcrumb() {
  const routerLocation = useLocation();
  const { pathname } = routerLocation;
  const { name } = useParams();

  if (pathname === '/' || pathname === '/admin/') {
    return <React.Fragment />;
  }

  // TODO: Get list of meshes, or mesh based on meshname and namespace

  let pathSnippets = pathname.split('/').filter(i => i);
  let breadcrumbItems = pathSnippets
    .map((_, index) => {
      const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;

      if (skippableNames[url]) {
        return undefined;
      } else if (breadcrumbLinkBackNames[url] !== undefined) {
        return (
          <AntdBreadcrumb.Item key={url}>
            <CrumbLink>
              <Link to={`${url}/`}>{breadcrumbLinkBackNames[url]}</Link>
            </CrumbLink>
          </AntdBreadcrumb.Item>
        );
      } else if (breakPathNames[pathSnippets[index]] !== undefined) {
        return (
          <AntdBreadcrumb.Item key={url}>
            <CrumbLink>
              <Link to={`${pathSnippets[index]}/`}>
                {breakPathNames[pathSnippets[index]]}
              </Link>
            </CrumbLink>
          </AntdBreadcrumb.Item>
        );
      } else if (breadcrumbNameMap[pathSnippets[index]] !== undefined) {
        return (
          <AntdBreadcrumb.Item key={url}>
            <CrumbLink>
              <Link to={`${url}/`}>
                {breadcrumbNameMap[pathSnippets[index]]}
              </Link>
            </CrumbLink>
          </AntdBreadcrumb.Item>
        );
      } else if (nonLinkNames[pathSnippets[index]] !== undefined) {
        return (
          <AntdBreadcrumb.Item key={url}>
            <CrumbLink>
              <span style={{ marginRight: '5px' }}>
                {nonLinkNames[pathSnippets[index]]}
              </span>{' '}
            </CrumbLink>
          </AntdBreadcrumb.Item>
        );
      } else {
        return (
          <AntdBreadcrumb.Item key={url}>
            <CrumbLink>
              <Link to={`${url}/`}>{pathSnippets[index]}</Link>
            </CrumbLink>
          </AntdBreadcrumb.Item>
        );
      }
    })
    .filter(crumb => crumb !== undefined);

  // do not make a breadcrumb for the namespace parts of the url
  for (let i = 0; i < breadcrumbItems.length; i++) {
    const endCrumb = (breadcrumbItems[i]!.key as string).split('/').pop()!;

    if (
      pathsHaveClutterAfter[endCrumb] !== undefined &&
      breadcrumbItems.length >= i + pathsHaveClutterAfter[endCrumb][0] // The least. We should be able to check any spot really
    ) {
      for (let j = pathsHaveClutterAfter[endCrumb].length - 1; j >= 0; j--) {
        breadcrumbItems.splice(i + pathsHaveClutterAfter[endCrumb][j], 1);
      }
    }
  }

  return (
    <Container isFloating={routerLocation.pathname === '/extensions/'}>
      <BreadcrumbContainer separator='>'>{breadcrumbItems}</BreadcrumbContainer>
    </Container>
  );
}
