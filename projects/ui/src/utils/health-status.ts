import { UpstreamStatus } from 'proto/github.com/solo-io/solo-apis/api/gloo/gloo/v1/upstream_pb';
import { FailoverSchemeStatus } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/v1/failover_pb';
import { PlacementStatus } from 'proto/github.com/solo-io/solo-projects/projects/gloo-fed/api/fed/core/v1/placement_pb';
import { colors } from 'Styles/colors';
import { Issue, NotificationType } from 'Components/Common/NotificationBox';

export enum StatusType {
  DEFAULT = 0,
  PLACEMENT = 1,
  FAILOVER = 2,
}

const SUCCESS_COLOR = { backgroundColor: colors.forestGreen };
const WARNING_COLOR = { backgroundColor: colors.sunGold };
const ERROR_COLOR = { backgroundColor: colors.grapefruitOrange };
const UNKNOWN_COLOR = {
  backgroundColor: 'transparent',
  borderColor: colors.juneGrey,
};

export const getHealthColor = (
  state?: number,
  statusType?: StatusType
): { backgroundColor: string; borderColor?: string } => {
  switch (statusType) {
    case StatusType.PLACEMENT: {
      switch (state) {
        case PlacementStatus.State.PLACED:
          return SUCCESS_COLOR;
        case PlacementStatus.State.INVALID:
        case PlacementStatus.State.FAILED:
          return ERROR_COLOR;
        case PlacementStatus.State.STALE:
          return WARNING_COLOR;
        case PlacementStatus.State.PENDING:
        case PlacementStatus.State.UNKNOWN:
          return UNKNOWN_COLOR;
      }
      break;
    }
    case StatusType.FAILOVER: {
      switch (state) {
        case FailoverSchemeStatus.State.ACCEPTED:
          return SUCCESS_COLOR;
        case FailoverSchemeStatus.State.INVALID:
        case FailoverSchemeStatus.State.FAILED:
          return ERROR_COLOR;
        case FailoverSchemeStatus.State.PENDING:
          return UNKNOWN_COLOR;
        case FailoverSchemeStatus.State.PROCESSING:
          return WARNING_COLOR;
      }
      break;
    }
    default: {
      switch (state) {
        case UpstreamStatus.State.ACCEPTED:
          return SUCCESS_COLOR;
        case UpstreamStatus.State.REJECTED:
          return ERROR_COLOR;
        case UpstreamStatus.State.PENDING:
          return UNKNOWN_COLOR;
        case UpstreamStatus.State.WARNING:
          return WARNING_COLOR;
      }
    }
  }
  return WARNING_COLOR;
};

const getNotificationType = (
  state?: number,
  statusType?: StatusType
): NotificationType | undefined => {
  const color = getHealthColor(state, statusType);
  if (color === WARNING_COLOR) {
    return NotificationType.WARNING;
  }
  if (color === ERROR_COLOR) {
    return NotificationType.ERROR;
  }
  return undefined;
};

type IssueWithType = {
  type: NotificationType;
  issue: Issue;
};

export const getNotificationIssue = (
  state?: number,
  statusType?: StatusType,
  reason?: string
): IssueWithType | undefined => {
  if (!reason) {
    return undefined;
  }
  const notificationType = getNotificationType(state, statusType);
  if (notificationType !== undefined) {
    return {
      type: notificationType,
      issue: { message: reason },
    };
  }
  return undefined;
};
