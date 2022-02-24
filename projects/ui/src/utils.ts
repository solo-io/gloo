import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Time } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';
import { format } from 'timeago.js';

export const formatTimestamp = (timestamp: Time.AsObject) => {
  const { nanos, seconds } = timestamp;
  if (!timestamp) {
    return '';
  }

  let dateObj = new Date(nanos + seconds * 1000);
  return format(dateObj, 'en_US');
};

export function copyTextToClipboard(copyText: string): Promise<void> {
  return navigator.clipboard.writeText(copyText);
}

// TODO - This is rough for now. See Issue 1207
export function getTimeAsSecondsString(fullTime?: Duration.AsObject): string {
  if (!fullTime) {
    return '-s';
  }

  // If there are partial seconds, round off to miliseconds
  if (fullTime.nanos) {
    if (!fullTime.seconds) {
      return `${fullTime.nanos}ns`;
    } else if (fullTime.nanos / 1000000 >= 0.51) {
      return `${fullTime.seconds}.${
        (fullTime.nanos / 1000000000)
          .toString()
          .match(/^-?\d+(?:\.\d{0,3})?/)![0]
      }s`;
    }

    return `${fullTime.seconds}s`;
  }

  return `${fullTime.seconds}s`;
}
