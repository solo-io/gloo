import React from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

type CountHealthProps = {
  healthy: boolean;
};
type CountProps = {
  count: number;
};
type CountBoxProps = {
  message: React.ReactNode;
} & CountHealthProps &
  CountProps;

const CountContainer = styled.div<CountHealthProps>`
  display: flex;
  align-items: center;
  border-radius: 6px;
  padding: 8px 15px;

  ${(props: CountHealthProps) =>
    props.healthy
      ? `
  color: ${colors.lakeBlue};
  background: ${colors.splashBlue};
  border: 1px solid ${colors.lakeBlue};`
      : `
  color: ${colors.pumpkinOrange};
  background: ${colors.tangerineOrange};
  border: 1px solid ${colors.grapefruitOrange};`}
`;
const CountNumber = styled.div<CountProps>`
  font-weight: 600;
  font-size: 24px;
  line-height: 24px;
  margin-right: 10px;

  min-width: 30px;
  ${(props: CountProps) =>
    `width: ${Math.floor(Math.log10(props.count) + 1) * 16}px`}
`;
const CountDescription = styled.div`
  flex: 1;
`;

export const CountBox = ({ count, message, healthy }: CountBoxProps) => {
  return (
    <CountContainer healthy={healthy}>
      <CountNumber count={count}>
        {count < 10 && '0'}
        {count}
      </CountNumber>
      <CountDescription>{message}</CountDescription>
    </CountContainer>
  );
};
