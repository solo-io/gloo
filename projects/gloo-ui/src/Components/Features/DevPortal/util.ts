export const secondsToString = (seconds: number | undefined): string => {
  if (!seconds) {
    return '';
  }
  const d = new Date(1970, 0, 1);
  d.setSeconds(seconds);
  return d.toLocaleString();
};
