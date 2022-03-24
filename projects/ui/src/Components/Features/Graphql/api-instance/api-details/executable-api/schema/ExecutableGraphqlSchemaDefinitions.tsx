import { Global } from '@emotion/core';
import { Collapse } from 'antd';
import { useGetGraphqlApiDetails } from 'API/hooks';
import { ReactComponent as GraphQLIcon } from 'assets/graphql-icon.svg';
import { SoloInput } from 'Components/Common/SoloInput';
import SoloNoMatches from 'Components/Common/SoloNoMatches';
import {
  DocumentNode,
  EnumTypeDefinitionNode,
  Kind,
  ObjectTypeDefinitionNode,
} from 'graphql';
import gql from 'graphql-tag';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect, useMemo, useState } from 'react';
import { makeSchemaDefinitionId } from 'utils/graphql-helpers';
import { ExeGqlEnumDefinition } from './enum-definition/ExeGqlEnumDefinition';
import { globalStyles } from './ExecutableGraphqlSchemaDefinitions.style';
import { ExeGqlObjectDefinition } from './object-definition/ExeGqlObjectDefinition';

const ExecutableGraphqlSchemaDefinitions: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { data: graphqlApi } = useGetGraphqlApiDetails(apiRef);

  // --- SCHEMA DEFINITIONS --- //
  type supportedDefinitionTypes =
    | ObjectTypeDefinitionNode
    | EnumTypeDefinitionNode;
  const schemaDefinitions = useMemo<supportedDefinitionTypes[]>(() => {
    // Try to parse the serialized GraphQL schema definition to JSON (using gql`...`).
    let query: DocumentNode;
    try {
      query = gql`
        ${graphqlApi?.spec?.executableSchema?.schemaDefinition}
      `;
    } catch {
      return [] as supportedDefinitionTypes[];
    }
    if (!query) return [] as supportedDefinitionTypes[];
    // We support enum and object type definitions here.
    const definitions = lodash.cloneDeep(
      query.definitions.filter(
        d =>
          d.kind === Kind.ENUM_TYPE_DEFINITION ||
          d.kind === Kind.OBJECT_TYPE_DEFINITION
      )
    ) as supportedDefinitionTypes[];
    // ? Uncomment this push(...mockEnumDefinitions) line for testing enums:
    // definitions.push(...mockEnumDefinitions);
    // We can sort the definitions here, and filtering will keep it sorted.
    definitions.sort((a, b) => {
      // Ordering: Query, mutation, Everything else.
      if (a.name.value === 'Query') return -1;
      else if (b.name.value === 'Query') return 1;
      if (a.name.value === 'Mutation') return -1;
      else if (b.name.value === 'Mutation') return 1;
      else return 0;
    });
    return definitions;
  }, [graphqlApi]);

  // --- SEARCH LOGIC --- //
  const [searchText, setSearchText] = useState('');
  const [filteredSchemaDefinitions, setFilteredSchemaDefinitions] = useState<
    typeof schemaDefinitions
  >([]);
  useEffect(() => {
    const lstext = searchText.toLowerCase();
    const newSchemaDefinitions = [] as typeof schemaDefinitions;
    // Check for special search cases:
    if (/type:.*/g.test(lstext)) {
      const processedText = searchText.replaceAll('type:', '');
      schemaDefinitions.forEach(d => {
        if (d.name.value === processedText) newSchemaDefinitions.push(d);
      });
    } else {
      schemaDefinitions.forEach(d => {
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
  }, [searchText, schemaDefinitions]);

  // --- COLLAPSE/ACCORDION PANEL LOGIC --- //
  const [openPanels, setOpenPanels] = useState<string | string[]>([]);
  useEffect(() => {
    if (filteredSchemaDefinitions.length > 0) {
      const idToOpen = makeSchemaDefinitionId(
        apiRef,
        filteredSchemaDefinitions[0] as supportedDefinitionTypes
      );
      setOpenPanels(idToOpen);
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
                    resolverType={d.name.value}
                    onReturnTypeClicked={t => setSearchText(`type:${t}`)}
                    schemaDefinitions={schemaDefinitions}
                    fields={d.fields ?? []}
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
