import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { format } from 'timeago.js';

export const formatTimestamp = (seconds: number | undefined) => {
  if (seconds === undefined) {
    return '';
  }

  let dateObj = new Date(seconds * 1000);
  return format(dateObj, 'en_US');
};

export function copyTextToClipboard(copyText: string): boolean {
  let textArea = document.createElement('textarea');

  textArea.style.position = 'fixed';
  textArea.style.top = '-999px';
  textArea.style.left = '-999px';

  // Ensure it has a small width and height. Setting to 1px / 1em
  // doesn't work as this gives a negative w/h on some browsers.
  textArea.style.width = '2em';
  textArea.style.height = '2em';

  // Avoid flash of white box if rendered for any reason.
  textArea.style.background = 'rgba(255, 255, 255, 0)';

  textArea.value = copyText;

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  let success = false;
  try {
    success = document.execCommand('copy');
  } catch (err) {
    console.log('Oops, unable to copy.' + err);
  }

  document.body.removeChild(textArea);

  return success;
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
