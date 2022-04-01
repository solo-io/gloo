import { CloseOutlined } from '@ant-design/icons';
import { Alert, Collapse } from 'antd';
import { useGetGraphqlApiDetails } from 'API/hooks';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import SoloAddButton from 'Components/Common/SoloAddButton';
import { SoloDropdown } from 'Components/Common/SoloDropdown';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import { StitchedSchema } from 'proto/github.com/solo-io/solo-apis/api/gloo/graphql.gloo/v1alpha1/graphql_pb';
import React, { useEffect, useMemo, useState } from 'react';
import { SoloNegativeButton } from 'Styles/StyledComponents/button';
import {
  isExecutableAPI,
  objectToArrayMap,
  parseSchemaDefinition,
} from 'utils/graphql-helpers';
import YAML from 'yaml';
import StitchedGqlAddSubGraphTypeMergeMapConfigItem from './StitchedGqlAddSubGraphTypeMergeMapConfigItem';

const sampleTypeMerge = `argsMap:
queryName:
selectionSet:`;

const StitchedGqlAddSubGraphTypeMergeMapConfig: React.FC<{
  onIsValidChange(isValid: boolean): void;
  onTypeMergeMapChange(
    typeMergeMap: [
      string,
      StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject
    ][]
  ): void;
  subGraphRef: ClusterObjectRef.AsObject;
}> = ({ onIsValidChange, subGraphRef, onTypeMergeMapChange }) => {
  const { data: subGraphqlApi } = useGetGraphqlApiDetails(subGraphRef);

  // --- TYPE MERGE MAP --- //
  const [typeMergeMap, setTypeMergeMap] = useState<
    {
      typeName: string;
      typeMergeConfig: string;
    }[]
  >([]);
  const [warningMessage, setWarningMessage] = useState('');
  useEffect(() => {
    // If there is a warning, we shouldn't be able to submit.
    onIsValidChange(warningMessage === '');
  }, [warningMessage]);

  useEffect(() => {
    // When types change, this resets the type merge map.
    // TODO: There should be a confirmation modal for this if the type merge config was edited.
    if (typeMergeMap.length === 0) return;
    setTypeMergeMap([]);
  }, [subGraphRef]);
  useEffect(() => {
    if (typeMergeMap.length === 0) return;
    // -- Parsing
    let parsedMap = [] as [
      string,
      StitchedSchema.SubschemaConfig.TypeMergeConfig.AsObject
    ][];
    for (let i = 0; i < typeMergeMap.length; i++) {
      const { typeName, typeMergeConfig } = typeMergeMap[i];
      let parsedMergeConfig: any;
      try {
        parsedMergeConfig = YAML.parse(typeMergeConfig);
        if (parsedMergeConfig.argsMap)
          parsedMergeConfig.argsMap = objectToArrayMap(
            parsedMergeConfig.argsMap
          );
        parsedMap.push([typeName, parsedMergeConfig]);
      } catch (err) {
        setWarningMessage(`${typeName}: ${(err as any).message}`);
        return;
      }
    }
    // -- Validation
    try {
      parsedMap.forEach(m => {
        const parsedMergeConfig = m[1];
        const configKeys = Object.keys(parsedMergeConfig);
        if (
          configKeys.length !== 3 ||
          !configKeys.includes('argsMap') ||
          !configKeys.includes('queryName') ||
          !configKeys.includes('selectionSet')
        )
          throw new Error(
            `${m[0]}): Must include values for 'argsMap', 'queryName', and 'selectionSet' only.`
          );
        // - argsMap
        if (parsedMergeConfig.argsMap === null) parsedMergeConfig.argsMap = [];
        else if (!parsedMergeConfig.argsMap.indexOf)
          throw new Error(`${m[0]}: Must include a valid 'argsMap'.`);
        // - queryName
        if (typeof parsedMergeConfig.queryName !== 'string')
          throw new Error(`${m[0]}: Must include a valid 'queryName'.`);
        // - selectionSet
        if (typeof parsedMergeConfig.selectionSet !== 'string')
          throw new Error(`${m[0]}: Must include a valid 'selectionSet'.`);
      });
      setWarningMessage('');
    } catch (err: any) {
      setWarningMessage(err.message);
    }
    onTypeMergeMapChange(parsedMap);
  }, [typeMergeMap]);

  // --- REMOVE TYPE MERGE MAPPING --- //
  const [confirmMapIdxToRemove, setConfirmMapIdxToRemove] = useState(-1);
  const removeFromTypeMergeMap = (index: number) => {
    const newTypeMergeMap = [...typeMergeMap];
    newTypeMergeMap.splice(index, 1);
    setTypeMergeMap(newTypeMergeMap);
  };

  // --- ADD TYPE MERGE MAPPING --- //
  // Gets the selected sub graph schema definition
  const subSchemaDefinitions = useMemo(() => {
    if (!subGraphqlApi) return [];
    if (isExecutableAPI(subGraphqlApi)) {
      setWarningMessage('');
      return parseSchemaDefinition(
        subGraphqlApi?.spec?.executableSchema?.schemaDefinition
      );
    } else {
      // TODO: This should work for stitched subgraphs as well (once the superschema is returned)
      setWarningMessage('Cannnot parse stitched schemas yet!');
      return [];
    }
  }, [subGraphqlApi]);
  // -- Sets up dropdown state
  const [newMergedTypeName, setNewMergedTypeName] = useState('');
  const availableTypes = useMemo(() => {
    // Create the type dropdown list from subschema definitions that have not been added yet.
    return subSchemaDefinitions
      .map(d => d.name.value)
      .filter(d => d !== 'Query' && !typeMergeMap.find(m => m.typeName === d));
  }, [subSchemaDefinitions, typeMergeMap]);
  useEffect(() => {
    if (availableTypes.length === 0) setNewMergedTypeName('');
    else setNewMergedTypeName(availableTypes[0]);
  }, [availableTypes]);

  // --- PANELS --- //
  const [openPanels, setOpenPanels] = useState<string[]>([]);

  return (
    <div className='block'>
      <div className='mt-5 mb-5 font-bold'>Type Merge Configuration</div>

      {availableTypes.length === 0 ? (
        <Alert
          type='success'
          showIcon
          className='mb-5'
          message={'All types added!'}
          description={' '}
        />
      ) : (
        <div className='mt-5 mb-5 flex'>
          <div className='flex items-center min-w-[200px]'>
            <div className='font-bold'>Type:&nbsp;&nbsp;</div>
            <SoloDropdown
              value={newMergedTypeName}
              options={availableTypes.map(t => ({
                key: t,
                value: t,
              }))}
              onChange={newValue => setNewMergedTypeName(newValue as string)}
              searchable={true}
            />
          </div>
          <div className='flex-grow flex justify-end items-center'>
            <SoloAddButton
              onClick={() => {
                const newTypeMergeMap = lodash.cloneDeep(typeMergeMap);
                newTypeMergeMap.push({
                  typeName: newMergedTypeName,
                  typeMergeConfig: sampleTypeMerge,
                });
                setOpenPanels([...openPanels, newMergedTypeName]);
                setTypeMergeMap(newTypeMergeMap);
              }}>
              Add Type Merge Configuration
            </SoloAddButton>
          </div>
        </div>
      )}

      {typeMergeMap.length > 0 && (
        <Collapse
          className='mt-5 mb-10'
          activeKey={openPanels}
          onChange={newOpenPanels => {
            if (typeof newOpenPanels === 'string')
              newOpenPanels = [newOpenPanels];
            // Any type in the availableTypes dropdown hasn't been added as a panel yet.
            // So we filter those out of the list (this removes deleted type merge mappings).
            newOpenPanels = newOpenPanels.filter(
              t => !availableTypes.includes(t)
            );
            setOpenPanels(newOpenPanels);
          }}>
          {typeMergeMap.map((m, idx) => (
            <Collapse.Panel
              key={m.typeName}
              header={
                <div className='grid grid-cols-[auto_min-content]'>
                  <div>{m.typeName}</div>
                  <div>
                    <SoloNegativeButton
                      style={{
                        width: 'auto',
                        minWidth: 'unset',
                        display: 'flex',
                        alignItems: 'center',
                        position: 'absolute',
                        top: '0px',
                        right: '0px',
                        bottom: '0px',
                        borderRadius: '0px',
                        padding: '0px 20px',
                      }}
                      onClick={e => {
                        e.stopPropagation();
                        // If the config was not changed or is empty, remove it.
                        // Otherwise, confirm.
                        const trimmedMergeConfig = m.typeMergeConfig.trim();
                        if (
                          trimmedMergeConfig === sampleTypeMerge ||
                          trimmedMergeConfig === ''
                        )
                          removeFromTypeMergeMap(confirmMapIdxToRemove);
                        else setConfirmMapIdxToRemove(idx);
                      }}>
                      <CloseOutlined />
                    </SoloNegativeButton>
                  </div>
                </div>
              }>
              <StitchedGqlAddSubGraphTypeMergeMapConfigItem
                schemaDefinitions={subSchemaDefinitions}
                typeMergeConfig={m.typeMergeConfig}
                onTypeMergeConfigChange={newValue => {
                  const newTypeMergeMap = [...typeMergeMap];
                  newTypeMergeMap[idx].typeMergeConfig = newValue;
                  setTypeMergeMap(newTypeMergeMap);
                }}
              />
            </Collapse.Panel>
          ))}
        </Collapse>
      )}

      {!!warningMessage && (
        <Alert
          showIcon
          // type='error'
          // message='Error'
          type='error'
          message='Error'
          // message={warningMessage}
          description={warningMessage}
        />
      )}

      <ConfirmationModal
        visible={confirmMapIdxToRemove !== -1}
        confirmPrompt='remove the edited type merge'
        confirmButtonText='Remove'
        goForIt={() => {
          removeFromTypeMergeMap(confirmMapIdxToRemove);
          setConfirmMapIdxToRemove(-1);
        }}
        cancel={() => setConfirmMapIdxToRemove(-1)}
        isNegative
      />
    </div>
  );
};

export default StitchedGqlAddSubGraphTypeMergeMapConfig;
