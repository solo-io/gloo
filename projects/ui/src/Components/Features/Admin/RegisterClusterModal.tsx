import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { ReactComponent as ClustersIcon } from 'assets/cluster-icon.svg';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import { ReactComponent as CopyIcon } from 'assets/document.svg';
import { SoloModal } from 'Components/Common/SoloModal';
import { colors } from 'Styles/colors';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { copyTextToClipboard } from 'utils';

const ModalContent = styled.div`
  padding: 25px 20px;
`;
const Title = styled.div`
  display: flex;
  font-size: 22px;
  line-height: 26px;
  font-weight: 500;
  margin-bottom: 20px;

  svg {
    margin-left: 8px;
  }
`;

const Hint = styled.div`
  background: ${colors.februaryGrey};
  border-radius: 8px;
  font-size: 14px;
  line-height: 44px;
  color: ${colors.septemberGrey};
  margin-bottom: 20px;
  padding: 0 11px;
`;

const InstructionsArea = styled.div`
  background: ${colors.splashBlue};
  border: 1px solid ${colors.seaBlue};
  border-radius: 8px;
  color: ${colors.seaBlue};
  padding: 15px 11px;
`;

const WarningCircle = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 16px;
  height: 16px;
  border-radius: 100%;
  background: ${colors.splashBlue};
  border: 1px solid ${colors.seaBlue};
  margin-right: 5px;

  svg {
    height: 10px !important;
    width: 3px;
    margin-right: 0 !important;

    * {
      fill: ${colors.seaBlue};
    }
  }
`;

const CommandArea = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  line-height: 60px;
  padding: 8px 11px;
  background: white;
  color: ${colors.novemberGrey};
  border-radius: 8px;
  margin: 10px 0;
`;

const FakeCommandLine = styled.div`
  display: flex;
  font-family: monospace;
`;

const FakeCommandLineNumber = styled.div`
  color: ${colors.juneGrey};
  margin-right: 8px;
`;

interface CopyIconHolderProps {
  copySuccessful: boolean | 'inactive';
}
const CopyIconHolder = styled.div<CopyIconHolderProps>`
  display: flex;
  height: 33px;
  width: 33px;
  justify-content: center;
  align-items: center;
  margin-left: 20px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.7s ease-out;

  ${(props: CopyIconHolderProps) => {
    if (props.copySuccessful === 'inactive') {
      return `background: ${colors.seaBlue};`;
    } else {
      return `background: ${
        props.copySuccessful ? colors.forestGreen : colors.pumpkinOrange
      };
      transition: background 0.2s ease-in;`;
    }
  }}

  svg {
    fill: white;
    height: 16px;
    margin-right: 0;
  }
`;

const InstructionHint = styled.div`
  font-size: 12px;
  line-height: 14px;
`;

interface Props {
  modalOpen: boolean;
  onClose: () => any;
}
export const RegisterClusterModal = ({ modalOpen, onClose }: Props) => {
  const [attemptedCopy, setAttemptedCopy] = useState<boolean | null>(null);

  const attemptCopyToClipboard = (copyText: string) => {
    setAttemptedCopy(copyTextToClipboard(copyText));

    setTimeout(() => {
      setAttemptedCopy(null);
    }, 500);
  };

  const registerCommand = 'glooctl cluster register';

  return (
    <SoloModal visible={modalOpen} width={625} onClose={onClose}>
      <ModalContent>
        <Title>
          Register a Kubernetes Cluster{' '}
          <IconHolder width={30}>
            <ClustersIcon />
          </IconHolder>
        </Title>
        <Hint>
          Check out the <a href='https://docs.solo.io/gloo/latest'>docs</a> for
          more info!
        </Hint>
        <InstructionsArea>
          <div>
            <WarningCircle>
              <WarningExclamation />
            </WarningCircle>{' '}
            Register a Cluster to Gloo Edge by running the following
            command:
          </div>

          <CommandArea>
            <FakeCommandLine>
              <FakeCommandLineNumber>01</FakeCommandLineNumber> ${' '}
              {registerCommand}
            </FakeCommandLine>
            <CopyIconHolder
              onClick={() => attemptCopyToClipboard(registerCommand)}
              copySuccessful={
                attemptedCopy === null ? 'inactive' : attemptedCopy
              }>
              <CopyIcon />
            </CopyIconHolder>
          </CommandArea>
          <InstructionHint>
            Be sure to download the latest version of glooctl.
          </InstructionHint>
        </InstructionsArea>
      </ModalContent>
    </SoloModal>
  );
};
