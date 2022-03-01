import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { Time } from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/common_pb';
import { format } from 'timeago.js';
import * as yup from 'yup';

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

// A regex matching a valid DNS-1123 compliant subdomain.
// Matched strings consist only of lower-case alphanumeric characters, '.', and '-',
// and must start and end with an alphanumeric character
// e.g. example.com        - matched
//      www.foo-bar.co.m  - matched
//      .aBc xyz-          - not matched
export const dnsCompliantSubdomainName =
  /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

export const nameValidationSchema = yup
  .string()
  .matches(
    /^[a-z0-9-.].*$/,
    "Only lowercase a-z, 0-9, '.', and '-' are allowed"
  )
  .matches(/^[a-z0-9](.*[a-z0-9])?$/, 'Must begin and end with a-z or 0-9')
  .matches(dnsCompliantSubdomainName, 'Must be DNS-1123 subdomain name')
  .max(253);
