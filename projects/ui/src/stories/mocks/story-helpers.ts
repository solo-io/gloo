import { jest } from '@storybook/jest';
import { injectable } from 'react-magnetic-di';
import { SWRResponse } from 'swr';

export function createLoadedSwrResponse<T>(data: T): SWRResponse<T> {
  return {
    data,
    error: undefined,
    isValidating: false,
    mutate: jest.fn(),
  };
}

export const createSwrInjectable = <ReturnType>(
  from: (...args: any[]) => SWRResponse<ReturnType>,
  implementation: ReturnType | (() => ReturnType)
) => {
  return injectable(from as any, () => {
    return createLoadedSwrResponse(
      typeof implementation === 'function'
        ? (implementation as any)()
        : implementation
    );
  });
};
