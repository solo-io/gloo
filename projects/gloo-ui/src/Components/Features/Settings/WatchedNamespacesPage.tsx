import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';

interface Props {}

export const WatchedNamespacesPage = (props: Props) => {
  const [watchedNamespacesList, setWatchedNamespacesList] = React.useState<
    string[]
  >(['fan', 'fun', 'jog']);

  const addNamespace = (newNamespace: string) => {
    const newArray = [...watchedNamespacesList];
    newArray.push(newNamespace);

    setWatchedNamespacesList(newArray);
  };

  const removeNamespace = (removeIndex: number) => {
    setWatchedNamespacesList([...watchedNamespacesList].splice(removeIndex, 1));
  };

  return (
    <SectionCard cardName={'Watched Namespaces'} logoIcon={<RelatedCircles />}>
      <StringCardsList
        values={watchedNamespacesList}
        valueDeleted={removeNamespace}
        createNew={addNamespace}
        createNewPromptText={'new.domain.com'}
      />
    </SectionCard>
  );
};
