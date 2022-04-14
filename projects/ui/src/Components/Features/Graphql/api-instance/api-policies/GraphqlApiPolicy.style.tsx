import { SoloButtonStyledComponent } from 'Styles/StyledComponents/button';
import styled from '@emotion/styled';

export const ItemWrapper = styled.div`
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: white;
  flex-direction: row;
  align-items: center;
  max-width: 750px;
`;

export const LabelWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

export const InputContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-self: center;
  margin-left: 50px;
  margin-top: 10px;
`;

export const InputLabel = styled.div`
  font-weight: bold;
`;

export const SwitchContainer = styled(InputContainer)`
  margin-left: 130px;
`;

export const NumericContainer = styled(InputContainer)`
  margin-left: 200px;
`;

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: flex-start;
  margin-top: 50px;
`;

export const DescriptionText = styled.div``;
