import React, { useEffect, useRef } from 'react';

export function useTimeout(callbackFn: () => any, delay: number) {
  const savedCallback = useRef();

  // Remember the latest callback.
  useEffect(() => {
    // @ts-ignore
    savedCallback.current = callbackFn;
  }, [callbackFn]);

  // Set up the interval.
  useEffect(() => {
    function func() {
      // @ts-ignore
      savedCallback.current();
    }
    if (delay !== null) {
      // @ts-ignore
      let timeoutId = setTimeout(func, delay);

      return () => clearTimeout(timeoutId);
    }
  }, [delay]);
}
