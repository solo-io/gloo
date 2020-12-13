import { colors } from './colors';
import { Status } from 'proto/github.com/solo-io/solo-kit/api/v1/status_pb';

export const healthConstants = {
  Good: {
    value: 1 as Status.StateMap[keyof Status.StateMap],
    color: colors.forestGreen,
    text: 'Accepted'
  },
  Pending: {
    value: 0 as Status.StateMap[keyof Status.StateMap],
    color: colors.sunGold,
    text: 'Pending'
  },
  Error: {
    value: 2 as Status.StateMap[keyof Status.StateMap],
    color: colors.grapefruitOrange,
    text: 'Rejected'
  }
};
