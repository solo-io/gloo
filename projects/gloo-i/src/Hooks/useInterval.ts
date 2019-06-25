import React, { useEffect, useRef } from 'react';

export function useInterval(callbackFn: () => any, delay: number) {
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
      let timeoutId = setInterval(func, delay);

      return () => clearTimeout(timeoutId);
    }
  }, [delay]);
}
