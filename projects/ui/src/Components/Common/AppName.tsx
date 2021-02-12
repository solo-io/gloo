import { useListGlooInstances } from 'API/hooks';
import React from 'react';
import { useLocation, Link, LinkProps } from 'react-router-dom';

export const AppName = () => {
  const { data: glooInstances, error: instancesError } = useListGlooInstances();

  return (
    <>{glooInstances && glooInstances.length > 1 ? 'Gloo Edge' : 'Gloo Fed'}</>
  );
};
