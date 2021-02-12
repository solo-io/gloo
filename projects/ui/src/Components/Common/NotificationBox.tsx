import React, { useState } from 'react';
import styled from '@emotion/styled/macro';
import { colors } from 'Styles/colors';
import { ReactComponent as ArrowDown } from 'assets/arrow-toggle.svg';
import { ReactComponent as ErrorIcon } from 'assets/big-unsuccessful-x.svg';
import { SoloLink } from 'Components/Common/SoloLink';
import { IconHolder } from 'Styles/StyledComponents/icons';
import { VerticalRule, HorizontalRule } from 'Styles/StyledComponents/shapes';

// overall container
type ColorProps = {
  color: string;
  backgroundColor: string;
  borderColor: string;
};
const Container = styled.div<ColorProps>`
  border-radius: 6px;
  padding: 8px 15px;
  margin-bottom: 15px;
  font-size: 14px;

  ${(props: ColorProps) =>
    `
  color: ${props.color};
  background: ${props.backgroundColor};
  border: 1px solid ${props.borderColor};`}
`;

// header section
const IssueHeaderContainer = styled.div`
  display: flex;
  align-items: center;
`;
const IssueMessageContainer = styled.div`
  margin-left: 8px;
`;
const IssueTypeContainer = styled.span`
  font-weight: bold;
`;
const ExpandLink = styled.div`
  cursor: pointer;
  font-weight: bold;
  flex: 1;
`;
const Toggler = styled.div`
  cursor: pointer;
`;
type ArrowIconProps = {
  expanded: boolean;
};
const ArrowIconHolder = styled(IconHolder)<ArrowIconProps>`
  ${(props: ArrowIconProps) => props.expanded && `transform: rotate(180deg);`}
`;

// expanded details section
const IssueDetailsContainer = styled.div`
  max-height: 150px;
  overflow-y: auto;
`;
const IssueContainer = styled.div`
  display: flex;
`;

// helpers
const getColors = (type: NotificationType): ColorProps => {
  switch (type) {
    case NotificationType.ERROR:
      return {
        color: colors.pumpkinOrange,
        backgroundColor: colors.tangerineOrange,
        borderColor: colors.grapefruitOrange,
      };
    case NotificationType.WARNING:
      return {
        color: colors.novaGold,
        backgroundColor: colors.flashlightGold,
        borderColor: colors.sunGold,
      };
    case NotificationType.NOTIFICATION:
    default:
      return {
        color: colors.lakeBlue,
        backgroundColor: colors.splashBlue,
        borderColor: colors.lakeBlue,
      };
  }
};

export enum NotificationType {
  NOTIFICATION = 0,
  WARNING = 1,
  ERROR = 2,
}

const notificationTypeStrings = {
  [NotificationType.NOTIFICATION]: 'Notification',
  [NotificationType.WARNING]: 'Warning',
  [NotificationType.ERROR]: 'Error',
};

export type Issue = {
  message: string;
  linkTitle?: string;
  detailsLink?: string;
};

// main component
type Props = {
  type: NotificationType;
  issues: Issue[];
  multipleIssuesMessage?: string;
};

export const NotificationBox = ({
  type,
  issues,
  multipleIssuesMessage = 'There are multiple issues.',
}: Props) => {
  const [isExpanded, setExpanded] = useState(false);

  const toggleExpanded = () => setExpanded(exp => !exp);

  const colors: ColorProps = getColors(type);
  const message =
    issues.length === 0
      ? ''
      : issues.length === 1
      ? issues[0].message
      : multipleIssuesMessage;
  const typeString = notificationTypeStrings[type];
  const headerLink =
    issues.length === 1 && issues[0].detailsLink ? (
      <SoloLink
        displayElement={issues[0].linkTitle ?? 'View Now'}
        link={issues[0].detailsLink}
        stylingOverrides={`font-weight: bold; color: ${colors.color};`}
      />
    ) : issues.length > 1 ? (
      <ExpandLink onClick={toggleExpanded}>
        {`${isExpanded ? 'Collapse' : 'Expand'} ${typeString} List`}
      </ExpandLink>
    ) : null;

  return (
    <Container
      color={colors.color}
      backgroundColor={colors.backgroundColor}
      borderColor={colors.borderColor}>
      <IssueHeaderContainer>
        <IconHolder
          width={20}
          applyColor={{ color: colors.color, strokeNotFill: true }}>
          <ErrorIcon />
        </IconHolder>
        <IssueMessageContainer>
          <IssueTypeContainer>{`${typeString}:`}</IssueTypeContainer> {message}
        </IssueMessageContainer>
        {headerLink ? (
          <>
            <VerticalRule color={colors.color} />
            {headerLink}
          </>
        ) : null}
        {issues.length > 1 && (
          <Toggler onClick={toggleExpanded}>
            <ArrowIconHolder
              expanded={isExpanded}
              width={15}
              applyColor={{ color: colors.color }}>
              <ArrowDown />
            </ArrowIconHolder>
          </Toggler>
        )}
      </IssueHeaderContainer>
      {isExpanded && (
        <>
          <HorizontalRule color={colors.borderColor} />
          <IssueDetailsContainer>
            <ul>
              {issues.map((issue, idx) => (
                <li key={`issue-${idx}`}>
                  <IssueContainer>
                    {issue.message}
                    {issue.detailsLink ? (
                      <>
                        <VerticalRule color={colors.color} />
                        <SoloLink
                          displayElement={issue.linkTitle ?? 'View Now'}
                          link={issue.detailsLink}
                          stylingOverrides={`font-weight: bold; color: ${colors.color};`}
                        />
                      </>
                    ) : null}
                  </IssueContainer>
                </li>
              ))}
            </ul>
          </IssueDetailsContainer>
        </>
      )}
    </Container>
  );
};
