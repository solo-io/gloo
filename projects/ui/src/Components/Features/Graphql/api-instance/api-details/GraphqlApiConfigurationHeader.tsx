import { graphqlConfigApi } from 'API/graphql';
import styled from '@emotion/styled/macro';
import { useGetGraphqlApiYaml } from 'API/hooks';
import AreaHeader from 'Components/Common/AreaHeader';
import { SoloModal } from 'Components/Common/SoloModal';
import VisualEditor from 'Components/Common/VisualEditor';
import YAML from 'yaml';
import { ClusterObjectRef } from 'proto/github.com/solo-io/skv2/api/core/v1/core_pb';
import React from 'react';
import {
  Footer,
  ModalContent,
  Title,
} from '../../new-api-modal/NewApiModal.style';
import {
  SoloButtonStyledComponent,
  SoloCancelButton,
  SoloNegativeButton,
} from 'Styles/StyledComponents/button';
import {
  UpdateGraphqlApiRequest,
  ValidateSchemaDefinitionRequest,
} from 'proto/github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/graphql_pb';
import ReactDiffViewer from 'react-diff-viewer';

const StyledFooter = styled(Footer)`
  margin-top: 20px;
  button {
    margin-right: 10px;
  }
`;
const GraphqlApiConfigurationHeader: React.FC<{
  apiRef: ClusterObjectRef.AsObject;
}> = ({ apiRef }) => {
  const { data: graphqlApiYaml } = useGetGraphqlApiYaml(apiRef);
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [graphqlSchema, setGraphqlSchema] = React.useState('');
  const [originalSchema, setOriginalSchema] = React.useState('');
  const [errorMessage, setErrorMessage] = React.useState('');
  const [showDiff, setShowDiff] = React.useState(false);

  const loadYaml = async () => {
    if (!apiRef.name || !apiRef.namespace) return '';
    try {
      const yaml = await graphqlConfigApi.getGraphqlApiYaml(apiRef);
      const parsedYaml = YAML.parse(yaml);
      const schemaDef = parsedYaml.spec.executableSchema.schemaDefinition;
      setGraphqlSchema(schemaDef);
      setOriginalSchema(schemaDef);
      return yaml;
    } catch (error) {
      console.error(error);
    }
    return '';
  };

  const onClose = () => {
    setShowEditModal(false);
  };

  const handleEditorConfig = () => {
    loadYaml().then(() => {
      setShowEditModal(true);
    });
  };

  const toggleShowDiff = () => {
    setShowDiff(!showDiff);
  };

  const handleReset = () => {
    setGraphqlSchema(originalSchema);
  };

  const handleClose = () => {
    setShowEditModal(false);
  };

  const handleSubmit = async () => {
    let schema = await graphqlConfigApi.getGraphqlApi(apiRef);
    schema.spec!.executableSchema!.schemaDefinition = graphqlSchema;

    const request = new UpdateGraphqlApiRequest();
    let validateRequest = new ValidateSchemaDefinitionRequest().toObject();
    const requestObj = request.toObject();
    requestObj.graphqlApiRef = apiRef;
    requestObj.spec = schema.spec;
    validateRequest = {
      ...validateRequest,
      ...requestObj,
    };

    await graphqlConfigApi
      .validateSchema(validateRequest)
      .then(() => {
        return graphqlConfigApi.updateGraphqlApi(requestObj).then(_res => {
          onClose();
        });
      })
      .catch(err => {
        setErrorMessage(err.message);
      });
  };
  /**
   * TODO:  Add in onEditConfig={handleEditorConfig}
   *        when you want to test graphql schema changes directly,
   *        intended mostly as a developer tool.
   */

  return (
    <div className='-mt-1 mb-5'>
      <AreaHeader
        title='Configuration'
        contentTitle={`${apiRef.namespace}--${apiRef.name}.yaml`}
        yaml={graphqlApiYaml}
        // onEditConfig={handleEditorConfig}
        onLoadContent={loadYaml}
      />
      <SoloModal visible={showEditModal} width={625} onClose={onClose}>
        <ModalContent>
          <Title>Edit Config</Title>
          {Boolean(errorMessage) && (
            <div className={`mb-5`}>
              <div className='p-2 text-orange-400 border border-orange-400 '>
                <div className='font-medium '>{errorMessage}</div>
              </div>
            </div>
          )}

          {!showDiff ? (
            <div>
              <VisualEditor
                theme='chrome'
                name='graphqlEditor'
                style={{
                  width: '100%',
                  maxHeight: '36vh',
                  cursor: 'text',
                }}
                onChange={(newValue, _e) => {
                  setGraphqlSchema(newValue);
                  // Change values.
                }}
                focus={true}
                fontSize={16}
                showPrintMargin={false}
                showGutter={true}
                highlightActiveLine={true}
                defaultValue={graphqlSchema || ''}
                value={graphqlSchema}
                readOnly={false}
                setOptions={{
                  highlightGutterLine: true,
                  showGutter: true,
                  fontFamily: 'monospace',
                  enableBasicAutocompletion: true,
                  enableLiveAutocompletion: true,
                  showLineNumbers: true,
                  tabSize: 2,
                }}
                mode='graphqlschema'
              />
            </div>
          ) : (
            <div>
              <ReactDiffViewer
                oldValue={originalSchema}
                splitView={true}
                newValue={graphqlSchema}
              />
            </div>
          )}
          <StyledFooter>
            <SoloButtonStyledComponent
              disabled={showDiff}
              onClick={handleSubmit as any}>
              Change Schema
            </SoloButtonStyledComponent>
            {originalSchema !== graphqlSchema && (
              <>
                <SoloCancelButton onClick={toggleShowDiff}>
                  {showDiff ? 'Hide Diff' : 'Show Diff'}
                </SoloCancelButton>
                <SoloCancelButton onClick={handleReset as any}>
                  Reset Schema
                </SoloCancelButton>
              </>
            )}
            <SoloNegativeButton onClick={handleClose as any}>
              Close
            </SoloNegativeButton>
          </StyledFooter>
        </ModalContent>
      </SoloModal>
    </div>
  );
};

export default GraphqlApiConfigurationHeader;
