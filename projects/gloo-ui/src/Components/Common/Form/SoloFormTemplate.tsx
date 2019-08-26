import styled from '@emotion/styled/macro';
import { Divider } from 'antd';
import * as React from 'react';

export const InputRow = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;
  padding: 10px 0 0;
`;

export const InputContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

export const Footer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
`;

const SectionHeader = styled.div`
  font-size: 18px;
  font-weight: 500;
  margin-top: 10px;
`;

const StyledDivider = styled(Divider)`
  margin: 12px 0;
`;

interface Props {
  formHeader?: string;
  children: React.ReactNode;
}

export const SoloFormTemplate: React.FC<Props> = props => {
  return (
    <div>
      {props.formHeader && (
        <React.Fragment>
          <SectionHeader>{props.formHeader}</SectionHeader>
          <StyledDivider />
        </React.Fragment>
      )}
      <InputContainer>{props.children}</InputContainer>
    </div>
  );
};
