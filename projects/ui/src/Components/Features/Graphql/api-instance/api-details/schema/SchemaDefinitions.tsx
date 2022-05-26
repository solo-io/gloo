import { Global } from '@emotion/core';
import { Collapse } from 'antd';
import { SoloInput } from 'Components/Common/SoloInput';
import SoloNoMatches from 'Components/Common/SoloNoMatches';
import { Kind } from 'graphql';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect, useState } from 'react';
import { Spacer } from 'Styles/StyledComponents/spacer';
import {
  getKindTypeReadableName,
  makeSchemaDefinitionId,
  SupportedDocumentNode,
} from 'utils/graphql-helpers';
import SchemaDefinitionContent from './SchemaDefinitionContent';
import { globalStyles } from './SchemaDefinitions.style';

const SchemaDefinitions: React.FC<{
  schema: SupportedDocumentNode;
  apiRef: ClusterObjectRef.AsObject;
  isEditable: boolean;
}> = ({ schema, apiRef, isEditable }) => {
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
        const name = d.name?.value ?? '';
        if (name === processedText) newSchemaDefinitions.push(d);
      });
    } else {
      schema.definitions.forEach(d => {
        const name = d.name?.value ?? '';
        if (name.toLowerCase().includes(lstext)) {
          newSchemaDefinitions.push(d);
        } else if (d.kind === Kind.ENUM_TYPE_DEFINITION) {
          const newValues = d.values?.filter(v =>
            v.name.value.toLowerCase().includes(lstext)
          );
          if (newValues && newValues.length > 0)
            newSchemaDefinitions.push({ ...d, values: newValues });
        } else if (d.kind === Kind.UNION_TYPE_DEFINITION) {
          const newTypes = d.types?.filter(v =>
            v.name.value.toLowerCase().includes(lstext)
          );
          if (newTypes && newTypes.length > 0)
            newSchemaDefinitions.push({ ...d, types: newTypes });
        } else if (d.kind === Kind.OPERATION_DEFINITION) {
          const newSelections = d.selectionSet.selections.filter(s => {
            if (
              s.kind === Kind.FIELD &&
              !s.name.value.toLowerCase().includes(lstext)
            )
              return false;
            return true;
          });
          if (newSelections && newSelections.length > 0)
            newSchemaDefinitions.push({
              ...d,
              selectionSet: {
                ...d.selectionSet,
                selections: newSelections,
              },
            });
        } else if (d.kind === Kind.DIRECTIVE_DEFINITION) {
          const newArgs = [] as any[];
          d.arguments?.forEach(a => {
            if (a.name.value.toLowerCase().includes(lstext)) newArgs.push(a);
          });
          if (newArgs.length > 0)
            newSchemaDefinitions.push({ ...d, arguments: newArgs });
        } else {
          const newFields = [] as any[];
          d.fields?.forEach(f => {
            if (f.name.value.toLowerCase().includes(lstext)) newFields.push(f);
          });
          if (newFields.length > 0)
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
                    <Spacer className='inline-block text-gray-600'>
                      {getKindTypeReadableName(d)}&nbsp;
                    </Spacer>
                    <div className='inline-block'>
                      {d.kind === Kind.DIRECTIVE_DEFINITION && '@'}
                      {d.name?.value ?? ''}
                    </div>
                  </div>
                }>
                <SchemaDefinitionContent
                  isEditable={isEditable}
                  schema={schema}
                  node={d}
                  onReturnTypeClicked={t => setSearchText(`type:${t}`)}
                />
              </Collapse.Panel>
            );
          })}
        </Collapse>
      )}
    </div>
  );
};

export default SchemaDefinitions;
