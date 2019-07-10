import { colors } from './colors';

export const healthConstants = {
  Good: {
    value: 1,
    color: colors.forestGreen,
    text: 'Accepted'
  },
  Pending: {
    value: 0,
    color: colors.sunGold,
    text: 'Pending'
  },
  Error: {
    value: 2,
    color: colors.grapefruitOrange,
    text: 'Rejected'
  }
};
