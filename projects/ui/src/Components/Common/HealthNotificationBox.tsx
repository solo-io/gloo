import React from 'react';
import { getNotificationIssue, StatusType } from 'utils/health-status';
import { NotificationBox } from 'Components/Common/NotificationBox';

type Props = {
  state?: number;
  statusType?: StatusType;
  reason?: string;
};

export const HealthNotificationBox = ({ state, statusType, reason }: Props) => {
  const notificationIssue = getNotificationIssue(state, statusType, reason);
  if (!notificationIssue) {
    return null;
  }
  return (
    <NotificationBox
      type={notificationIssue.type}
      issues={[notificationIssue.issue]}
    />
  );
};
