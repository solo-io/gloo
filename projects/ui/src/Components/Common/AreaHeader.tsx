import React, { useState } from 'react';
import styled from '@emotion/styled';
import { colors } from 'Styles/colors';
import { ReactComponent as DownloadIcon } from 'assets/document.svg';
import { ReactComponent as ViewIcon } from 'assets/view-icon.svg';
import { doDownload } from 'download-helper';
import YamlDisplayer from 'Components/Common/YamlDisplayer';
import { AreaTitle } from 'Styles/StyledComponents/headings';
import { HealthIndicator } from './HealthIndicator';
import { StatusType } from 'utils/health-status';
import { VerticalRule } from 'Styles/StyledComponents/shapes';

const RowContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 16px;
`;

const Actionables = styled.div`
  display: flex;
  align-items: center;
  color: ${colors.seaBlue};

  > div {
    display: flex;
    align-items: center;
    margin-left: 20px;
    cursor: pointer;
  }

  svg {
    margin-right: 5px;
  }
`;

const HealthTitle = styled.div`
  margin-right: 10px;
`;

type Props = {
  title: string;
  contentTitle?: string;
  yaml?: string;
  onLoadContent?: () => Promise<string>;
  health?: {
    state: number;
    type?: StatusType;
    title?: string;
    reason?: string;
  };
};

const AreaHeader = ({
  title,
  contentTitle = 'unnamed',
  onLoadContent,
  yaml,
  health,
}: Props) => {
  const [isExpanded, setExpanded] = useState(false);
  const [contentValue, setContentValue] = useState(yaml ?? '');

  const ensureContentLoaded = async () => {
    if (contentValue || !onLoadContent) {
      return contentValue;
    }
    const value = await onLoadContent();
    setContentValue(value);
    return value;
  };

  const onClickView = async () => {
    await ensureContentLoaded();
    setExpanded(expanded => !expanded);
  };

  const onClickDownload = async () => {
    const value = await ensureContentLoaded();
    doDownload(value, contentTitle);
  };

  React.useEffect(() => {
    if (yaml) {
      setContentValue(yaml);
    }
  }, [yaml]);
  return (
    <>
      <RowContainer>
        <AreaTitle style={{ flex: 1 }}>{title}</AreaTitle>
        {onLoadContent && (
          <Actionables>
            <div onClick={onClickView}>
              <ViewIcon /> {isExpanded ? 'Hide' : 'View'} Raw Config
            </div>
            <div onClick={onClickDownload}>
              <DownloadIcon /> {contentTitle}
            </div>
          </Actionables>
        )}
        {onLoadContent && health && <VerticalRule />}
        {health?.title && <HealthTitle>{health.title}</HealthTitle>}
        {health && (
          <HealthIndicator
            healthStatus={health.state}
            statusType={health.type}
            issueText={health.reason}
          />
        )}
      </RowContainer>
      {isExpanded ? (
        <div className='mb-5'>
          <YamlDisplayer contentString={contentValue} />
        </div>
      ) : null}
    </>
  );
};

export default AreaHeader;
