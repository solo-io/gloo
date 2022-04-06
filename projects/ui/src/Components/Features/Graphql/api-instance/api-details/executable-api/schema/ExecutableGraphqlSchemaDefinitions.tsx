import { Global } from '@emotion/core';
import { Collapse } from 'antd';
import { useGetGraphqlApiDetails } from 'API/hooks';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { SoloInput } from 'Components/Common/SoloInput';
import SoloNoMatches from 'Components/Common/SoloNoMatches';
import { Kind } from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect, useMemo, useState } from 'react';
import { getParsedSchema, makeSchemaDefinitionId } from 'utils/graphql-helpers';
import { ExeGqlEnumDefinition } from './enum-definition/ExeGqlEnumDefinition';
import { globalStyles } from './ExecutableGraphqlSchemaDefinitions.style';
import { ExeGqlObjectDefinition } from './object-definition/ExeGqlObjectDefinition';

const ExecutableGraphqlSchemaDefinitions: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  // --- SCHEMA DEFINITIONS --- //
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);
  const schema = useMemo(() => getParsedSchema(graphqlApi), [graphqlApi]);

  // --- SEARCH LOGIC --- //
  const [searchText, setSearchText] = useState('');
  const [filteredSchemaDefinitions, setFilteredSchemaDefinitions] = useState<
    typeof schema.definitions
  >([]);
  useEffect(() => {
    const lstext = searchText.toLowerCase();
    const newSchemaDefinitions = [] as typeof schema.definitions;
    // Check for special search cases:
    if (/type:.*/g.test(lstext)) {
      const processedText = searchText.replaceAll('type:', '');
      schema.definitions.forEach(d => {
        if (d.name.value === processedText) newSchemaDefinitions.push(d);
      });
    } else {
      schema.definitions.forEach(d => {
        if (d.name.value.toLowerCase().includes(lstext)) {
          newSchemaDefinitions.push(d);
        } else if (d.kind === Kind.ENUM_TYPE_DEFINITION) {
          const newValues = d.values?.filter(v =>
            v.name.value.toLowerCase().includes(lstext)
          );
          if (newValues && newValues.length > 0)
            newSchemaDefinitions.push({ ...d, values: newValues });
        } else {
          // d.kind === Kind.OBJECT_TYPE_DEFINITION
          const newFields = d.fields?.filter(v =>
            v.name.value.toLowerCase().includes(lstext)
          );
          if (newFields && newFields.length > 0)
            newSchemaDefinitions.push({ ...d, fields: newFields });
        }
      });
    }
    setFilteredSchemaDefinitions(newSchemaDefinitions);
  }, [searchText, schema]);

  // --- COLLAPSE/ACCORDION PANEL LOGIC --- //
  const [openPanels, setOpenPanels] = useState<string | string[]>([]);
  useEffect(() => {
    if (filteredSchemaDefinitions.length > 0) {
      const idToOpen = makeSchemaDefinitionId(
        apiRef,
        filteredSchemaDefinitions[0]
      );
      // setOpenPanels(idToOpen);
      if (!openPanels.includes(idToOpen))
        setOpenPanels([idToOpen, ...openPanels]);
    }
  }, [filteredSchemaDefinitions]);

  return (
    <div className='relative mb-5'>
      <Global styles={globalStyles} />

      <div className='max-w-[500px] mb-5'>
        <SoloInput
          placeholder='Filter by name...'
          value={searchText}
          onChange={s => setSearchText(s.target.value)}
        />
      </div>

      {filteredSchemaDefinitions.length === 0 ? (
        <>
          <div className='pt-1 pb-2'>
            <SoloNoMatches />
          </div>
          <hr />
        </>
      ) : (
        <Collapse
          activeKey={openPanels}
          onChange={newIds => setOpenPanels(newIds)}>
          {filteredSchemaDefinitions.map(d => {
            const definitionId = makeSchemaDefinitionId(apiRef, d);
            return (
              <Collapse.Panel
                key={definitionId}
                id={definitionId}
                header={
                  <div className='inline font-medium text-gray-900 whitespace-nowrap'>
                    <GraphQLIcon className='w-4 h-4 fill-current inline' />
                    &nbsp;&nbsp;
                    {d.name.value}
                  </div>
                }>
                {d.kind === Kind.ENUM_TYPE_DEFINITION ? (
                  <ExeGqlEnumDefinition
                    resolverType={d.name.value}
                    values={d.values ?? []}
                  />
                ) : (
                  <ExeGqlObjectDefinition
                    apiRef={apiRef}
                    schema={schema}
                    objectTypeDefinition={d}
                    onReturnTypeClicked={t => setSearchText(`type:${t}`)}
                  />
                )}
              </Collapse.Panel>
            );
          })}
        </Collapse>
      )}
    </div>
  );
};

export default ExecutableGraphqlSchemaDefinitions;
