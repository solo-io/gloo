import { CloseOutlined } from '@ant-design/icons';
import { Alert, Collapse } from 'antd';
import ConfirmationModal from 'Components/Common/ConfirmationModal';
import lodash from 'lodash';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React, { useEffect, useState } from 'react';
import { SoloNegativeButton } from 'Styles/StyledComponents/button';
import StitchedGqlTypeMergeFieldDropdown from './StitchedGqlTypeMergeFieldDropdown';
import StitchedGqlAddSubGraphTypeMergeMapConfigItem from './StitchedGqlTypeMergeMapConfigItem';
import {
  ParsedTypeMergeMap,
  typeMergeMapFromStringFormat,
  TypeMergeMapStringFormat,
  typeMergeMapToStringFormat,
  validateTypeMergeMap,
} from './StitchedGqlTypeMergeMapHelpers';

// TODO: Fix argsmap > args naming.
const sampleTypeMerge = `argsMap:
queryName:
selectionSet:`;

const StitchedGqlTypeMergeMapConfig: React.FC<{
  onIsValidChange(isValid: boolean): void;
  initialTypeMergeMap: ParsedTypeMergeMap;
  onTypeMergeMapChange(typeMergeMap: ParsedTypeMergeMap): void;
  subGraphqlApiRef: ClusterObjectRef.AsObject;
}> = ({
  onIsValidChange,
  initialTypeMergeMap,
  onTypeMergeMapChange,
  subGraphqlApiRef,
}) => {
  // --- TYPE MERGE MAP (SF = string formatted) --- //
  const [typeMergeMapSF, setTypeMergeMapSF] =
    useState<TypeMergeMapStringFormat>([]);
  useEffect(() => {
    setTypeMergeMapSF(typeMergeMapToStringFormat(initialTypeMergeMap));
  }, []);
  useEffect(() => {
    try {
      // Parse
      const parsedMap = typeMergeMapFromStringFormat(typeMergeMapSF);
      // Call event handlers
      onTypeMergeMapChange(parsedMap);
      // Validate (this can throw errors, which we handle here)
      validateTypeMergeMap(parsedMap);
      // Clear the warning
      setWarningMessage('');
    } catch (err: any) {
      setWarningMessage(err.message);
    }
  }, [typeMergeMapSF]);

  // --- REMOVE TYPE MERGE MAPPING --- //
  const [confirmMapIdxToRemove, setConfirmMapIdxToRemove] = useState(-1);
  const removeFromTypeMergeMap = (index: number) => {
    const newTypeMergeMap = [...typeMergeMapSF];
    newTypeMergeMap.splice(index, 1);
    setTypeMergeMapSF(newTypeMergeMap);
  };

  // --- WARNING MESSAGE --- //
  const [warningMessage, setWarningMessage] = useState('');
  useEffect(() => {
    // If there is a warning, we shouldn't be able to submit.
    onIsValidChange(warningMessage === '');
  }, [warningMessage]);

  // --- PANELS --- //
  const [openPanels, setOpenPanels] = useState<string[]>([]);

  return (
    <div className='block'>
      <div className='mt-5 mb-5 font-bold'>Type Merge Configuration</div>

      {/* --- FIELD DROPDOWN --- */}
      <StitchedGqlTypeMergeFieldDropdown
        subGraphqlApiRef={subGraphqlApiRef}
        addedTypeNames={typeMergeMapSF.map(m => m.typeName)}
        onAddTypeMerge={(newMergedTypeName: string) => {
          const newTypeMergeMap = lodash.cloneDeep(typeMergeMapSF);
          newTypeMergeMap.push({
            typeName: newMergedTypeName,
            typeMergeConfig: sampleTypeMerge,
          });
          setOpenPanels([...openPanels, newMergedTypeName]);
          setTypeMergeMapSF(newTypeMergeMap);
        }}
      />

      {/* --- TYPE MERGE CONFIGS --- */}
      {typeMergeMapSF.length > 0 && (
        <Collapse
          className='mt-5 mb-10'
          activeKey={openPanels}
          onChange={newOpenPanels => {
            if (typeof newOpenPanels === 'string')
              newOpenPanels = [newOpenPanels];
            // This removes any deleted type merge mappings.
            newOpenPanels = newOpenPanels.filter(
              t => typeMergeMapSF.find(m => m.typeName === t) !== undefined
            );
            setOpenPanels(newOpenPanels);
          }}>
          {typeMergeMapSF.map((m, idx) => (
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
                        const trimmedMergeConfig = m.typeMergeConfig.trim();
                        if (
                          trimmedMergeConfig === sampleTypeMerge ||
                          trimmedMergeConfig === ''
                        )
                          removeFromTypeMergeMap(confirmMapIdxToRemove);
                        // Otherwise, confirm removing it.
                        else setConfirmMapIdxToRemove(idx);
                      }}>
                      <CloseOutlined />
                    </SoloNegativeButton>
                  </div>
                </div>
              }>
              <div data-testid={`type-merge-${m.typeName}`}>
                <StitchedGqlAddSubGraphTypeMergeMapConfigItem
                  typeMergeConfig={m.typeMergeConfig}
                  onTypeMergeConfigChange={newValue => {
                    const newTypeMergeMap = [...typeMergeMapSF];
                    newTypeMergeMap[idx].typeMergeConfig = newValue;
                    setTypeMergeMapSF(newTypeMergeMap);
                  }}
                />
              </div>
            </Collapse.Panel>
          ))}
        </Collapse>
      )}

      {/* --- ALERTS + CONFIRMATION --- */}
      {!!warningMessage && (
        <Alert
          showIcon
          type='error'
          message='Error'
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

export default StitchedGqlTypeMergeMapConfig;
