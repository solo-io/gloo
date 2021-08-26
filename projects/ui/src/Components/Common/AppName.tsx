import React from 'react';
import { useIsGlooFedEnabled } from 'API/hooks';

export const AppName = () => {
  const {
    data: glooFedCheckResponse,
    error: glooFedCheckError,
  } = useIsGlooFedEnabled();

  if (glooFedCheckError) {
    console.error('Could not check if Gloo Fed is enabled', glooFedCheckError);
  }

  return <>{glooFedCheckResponse?.enabled ? 'Gloo Fed' : 'Gloo Edge'}</>;
};
