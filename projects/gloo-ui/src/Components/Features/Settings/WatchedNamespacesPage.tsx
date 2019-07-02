import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { colors } from 'Styles';
import { SectionCard } from 'Components/Common/SectionCard';
import { StringCardsList } from 'Components/Common/StringCardsList';
import { ReactComponent as RelatedCircles } from 'assets/related-circles.svg';
import { GetSettingsRequest } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/config_pb';
import { useGetSettings } from 'Api';

interface Props {}

export const WatchedNamespacesPage = (props: Props) => {
  let req = new GetSettingsRequest();
  const { data, loading, error } = useGetSettings(req);

  const [watchedNamespacesList, setWatchedNamespacesList] = React.useState<
    string[]
  >([]);

  React.useEffect(() => {
    if (data && data.settings) {
      setWatchedNamespacesList(data.settings.watchNamespacesList);
    }
  }, [loading]);

  if (!data || loading) {
    return <div>Loading...</div>;
  }

  const addNamespace = (newNamespace: string) => {
    const newArray = [...watchedNamespacesList];
    newArray.push(newNamespace);

    setWatchedNamespacesList(newArray);
  };

  const removeNamespace = (removeIndex: number) => {
    let newList = [...watchedNamespacesList];
    newList.splice(removeIndex, 1);
    setWatchedNamespacesList(newList);
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
