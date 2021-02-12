import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';

type VerticalRuleProps = {
  height?: string;
  color?: string;
};

export const VerticalRule = styled.div<VerticalRuleProps>`
  border-left: 1px solid
    ${(props: VerticalRuleProps) => props.color ?? colors.septemberGrey};
  margin: 0 10px;
  height: ${(props: VerticalRuleProps) => props.height ?? '1em'};
`;

type HorizontalRuleProps = {
  color?: string;
};

export const HorizontalRule = styled.div<HorizontalRuleProps>`
  height: 1px;
  width: 100%;
  background: ${(props: HorizontalRuleProps) =>
    props.color ?? colors.septemberGrey};
  margin: 10px 0;
`;
