import { Alert } from 'antd';
import { graphqlConfigApi } from 'API/graphql';
import { useGetGraphqlApiDetails } from 'API/hooks';
import SoloAddButton from 'Components/Common/SoloAddButton';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect, useMemo, useState } from 'react';
import {
  SupportedDefinitionNode,
  isExecutableAPI,
  getParsedExecutableApiSchema,
  parseSchemaString,
} from 'utils/graphql-helpers';

const StitchedGqlTypeMergeFieldDropdown: React.FC<{
  subGraphqlApiRef: ClusterObjectRef.AsObject;
  onAddTypeMerge(newMergedTypeName: string): void;
  addedTypeNames: string[];
}> = ({ subGraphqlApiRef, onAddTypeMerge, addedTypeNames }) => {
  const { data: subGraphqlApi } = useGetGraphqlApiDetails(subGraphqlApiRef);

  // --- GET SELECTED SCHEMA DEFINITIONS --- //
  // These are the options in the dropdown.
  const [subSchemaDefinitions, setSubSchemaDefinitions] =
    useState<SupportedDefinitionNode[]>();
  useEffect(() => {
    if (!subGraphqlApi) return;
    if (isExecutableAPI(subGraphqlApi)) {
      // For executable APIs, we have the schema from the useGetGraphqlDetails call.
      const newSubSchemaDefinitions =
        getParsedExecutableApiSchema(subGraphqlApi).definitions;
      setSubSchemaDefinitions(newSubSchemaDefinitions);
    } else {
      // For stitched APIs, we have to get the schema string separately.
      graphqlConfigApi
        .getStitchedSchemaDefinition(subGraphqlApiRef)
        .then(schemaString => {
          const newSubSchemaDefinitions =
            parseSchemaString(schemaString).definitions;
          setSubSchemaDefinitions(newSubSchemaDefinitions);
        });
    }
  }, [subGraphqlApi]);

  // --- DROPDOWN STATE --- //
  const [newMergedTypeName, setNewMergedTypeName] = useState('');
  // availableTypes are the object type definitions for this subschema that have not been added yet.
  const availableTypes = useMemo(() => {
    if (!subSchemaDefinitions) return [];
    return subSchemaDefinitions
      .map(d => d.name.value)
      .filter(
        d =>
          // "Query" and "Mutation" are special types:
          // https://graphql.org/learn/schema/#the-query-and-mutation-types
          d !== 'Query' && d !== 'Mutation' && !addedTypeNames.includes(d)
      );
  }, [subSchemaDefinitions, addedTypeNames]);
  useEffect(() => {
    if (availableTypes.length === 0) setNewMergedTypeName('');
    else setNewMergedTypeName(availableTypes[0]);
  }, [availableTypes]);

  if (!subGraphqlApi || subSchemaDefinitions === undefined) return null;
  if (availableTypes.length === 0)
    return (
      <Alert
        type='success'
        showIcon
        className='mb-5'
        message={'All types added!'}
        description={' '}
      />
    );
  return (
    <div className='mt-5 mb-5 flex'>
      <div className='flex items-center min-w-[200px]'>
        <div className='font-bold'>Type:&nbsp;&nbsp;</div>
        <SoloDropdown
          data-testid='type-merge-name-dropdown'
          value={newMergedTypeName}
          options={availableTypes.map(t => ({ key: t, value: t }))}
          onChange={newValue => setNewMergedTypeName(newValue as string)}
          searchable={true}
        />
      </div>
      <div className='flex-grow flex justify-end items-center'>
        <SoloAddButton
          data-testid='add-type-merge-button'
          onClick={() => onAddTypeMerge(newMergedTypeName)}>
          Add Type Merge Configuration
        </SoloAddButton>
      </div>
    </div>
  );
};

export default StitchedGqlTypeMergeFieldDropdown;
