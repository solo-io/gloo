import React, { useEffect, useRef } from 'react';

export function useInterval(callbackFn: () => any, delay: number | null) {
  const savedCallback = useRef<ReturnType<typeof callbackFn>>();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callbackFn;
  }, [callbackFn]);

  // Set up the interval.
  useEffect(() => {
    function func() {
      savedCallback.current();
    }
    if (delay !== null) {
      let timeoutId = setInterval(func, delay);

      return () => clearInterval(timeoutId);
    }
  }, [delay]);
}
