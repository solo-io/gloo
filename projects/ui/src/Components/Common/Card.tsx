import styled from '@emotion/styled';
import { colors } from 'Styles/colors';

export const Card = styled.div`
  position: relative;
  width: 100%;
  padding: 23px 18px 18px;
  background: white;
  border-radius: 10px;
  box-shadow: 0px 4px 9px ${colors.boxShadow};
  box-sizing: border-box;
`;

export const CardSubsectionWrapper = styled.div`
  border-radius: 7px;
  padding: 18px;
  background: ${colors.januaryGrey};
`;

export const CardSubsectionContent = styled.div`
  border-radius: 7px;
  padding: 18px 20px;
  background: white;
  height: 100%;
`;

export const CardWhiteSubsection = styled(CardSubsectionContent)`
  background: white;
  border: 1px solid ${colors.marchGrey};
`;

type CardSubsectionDividerProps = {
  vertical?: boolean;
};
export const CardSubsectionDivider = styled.div<CardSubsectionDividerProps>`
  border${(props: CardSubsectionDividerProps) =>
    props.vertical ? '-right: ' : '-top: '}1px solid ${colors.marchGrey};
  ${(props: CardSubsectionDividerProps) =>
    !props.vertical && 'margin-bottom: 20px;'}
`;

/* CARD HEADER BLOCKS */

export const CardHeader = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  background: ${colors.marchGrey};
  padding: 0 13px;
  line-height: 58px;
  border-radius: 10px 10px 0 0;
`;

export const CardFooter = styled.div`
  display: flex;
  border-radius: 0 0 10px 10px;
  background: ${colors.februaryGrey};
  color: ${colors.septemberGrey};
  font-size: 14px;
  line-height: 30px;
  padding-left: 11px;

  > div {
    margin-right: 8px;
  }
`;
