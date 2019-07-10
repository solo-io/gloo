import styled from '@emotion/styled/macro';
import { Divider } from 'antd';
import { InputContainer } from 'Components/Features/Upstream/Creation/CreateUpstreamForm';
import * as React from 'react';

const InputItem = styled.div`
  display: flex;
  flex-direction: column;
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
  children: React.ReactNode[];
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
      <InputContainer>
        {React.Children.map(props.children, field => (
          <InputItem>{field}</InputItem>
        ))}
      </InputContainer>
    </div>
  );
};
