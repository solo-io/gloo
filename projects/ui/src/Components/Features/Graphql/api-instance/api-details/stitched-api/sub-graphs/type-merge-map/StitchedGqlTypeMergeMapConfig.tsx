import { CloseOutlined } from '@ant-design/icons';
import { Collapse } from 'antd';
import { useConfirm } from 'Components/Context/ConfirmModalContext';
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
} from './StitchedGqlTypeMergeMapHelpers';

const defaultTypeMergeConfig = `args:
queryName:
selectionSet:`;

export const typeMergeConfigChangedFromDefault = (typeMergeConfig: string) => {
  return (
    typeMergeConfig.replaceAll(/\s|\n|\t/g, '') !==
    defaultTypeMergeConfig.replaceAll(/\s|\n|\t/g, '')
  );
};

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
  const confirm = useConfirm();

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
      // If it got to here it is valid
      onIsValidChange(true);
    } catch (_) {
      onIsValidChange(false);
    }
  }, [typeMergeMapSF]);

  // --- REMOVE TYPE MERGE MAPPING --- //
  const removeFromTypeMergeMap = (index: number) => {
    const newTypeMergeMap = [...typeMergeMapSF];
    newTypeMergeMap.splice(index, 1);
    setTypeMergeMapSF(newTypeMergeMap);
  };

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
            typeMergeConfig: defaultTypeMergeConfig,
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
                        if (
                          !typeMergeConfigChangedFromDefault(
                            m.typeMergeConfig
                          ) ||
                          m.typeMergeConfig.trim() === ''
                        )
                          removeFromTypeMergeMap(idx);
                        // Otherwise, confirm removing it.
                        else
                          confirm({
                            confirmPrompt: 'remove the edited type merge',
                            confirmButtonText: 'Remove',
                            isNegative: true,
                          }).then(() => removeFromTypeMergeMap(idx));
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
    </div>
  );
};

export default StitchedGqlTypeMergeMapConfig;
