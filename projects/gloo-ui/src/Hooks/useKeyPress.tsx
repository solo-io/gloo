import { useState, useEffect, EventHandler } from 'react';

export function useKeyPress(targetKey: string, onKeyUpEvent: () => void) {
  // State for keeping track of whether key is pressed
  const [keyPressed, setKeyPressed] = useState(false);
  // If pressed key is our target key then set to true
  const downHandler = ({ key }: any): any => {
    if (!!targetKey.length && key === targetKey) {
      setKeyPressed(true);
    }
  };

  // If released key is our target key then set to false
  const upHandler = ({ key }: any): any => {
    if (!!targetKey.length && key === targetKey) {
      onKeyUpEvent();
      setKeyPressed(false);
    }
  };

  // Add event listeners
  useEffect(() => {
    window.addEventListener('keydown', downHandler);
    window.addEventListener('keyup', upHandler);
    // Remove event listeners on cleanup
    return () => {
      window.removeEventListener('keydown', downHandler);
      window.removeEventListener('keyup', upHandler);
    };
  }, []); // Empty array ensures that effect is only run on mount and unmount

  return keyPressed;
}
