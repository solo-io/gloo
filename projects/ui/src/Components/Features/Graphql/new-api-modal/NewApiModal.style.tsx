import styled from '@emotion/styled';

export const ModalContent = styled.div`
  padding: 25px 20px;
`;

export const Title = styled.div`
  display: flex;
  font-size: 22px;
  line-height: 26px;
  font-weight: 500;
  margin-bottom: 20px;

  svg {
    margin-left: 8px;
  }
`;

export const InputWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding-bottom: 1.5rem;
`;

export const Footer = styled.footer`
  display: flex;
  flex-direction: row-reverse;
`;

export const StyledWarning = styled.div`
  backgroundcolor: #fef2f2;
`;
