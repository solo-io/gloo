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

const CrumbLink = styled.span`
  color: ${colors.seaBlue};
  font-size: 14px;
  cursor: pointer;
`;

export interface RouteParams {
  //... eg, virtualservice?: string
}

function BreadcrumbC({
  history,
  match,
  location
}: RouteComponentProps<RouteParams>) {
  return <BreadcrumbContainer separator='>' />;
}

export const Breadcrumb = withRouter(BreadcrumbC);
