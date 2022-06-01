import { SoloFormFileUpload } from 'Components/Common/SoloFormComponents';
import VisualEditor, {
  SoloFormVisualEditorProps,
} from 'Components/Common/VisualEditor';
import { useFormikContext } from 'formik';
import React from 'react';
import { Spacer } from 'Styles/StyledComponents/spacer';
import WarningMessage from '../../../executable-api/WarningMessage';
import { isBase64, base64ToString, stringToBase64 } from '../converters';
import { ResolverWizardFormProps } from '../ResolverWizard';

interface GrpcProtoCheckProps {
  setWarningMessage: (value: string) => any;
  warningMessage: string;
}

export const GrpcProtoCheck = (props: GrpcProtoCheckProps) => {
  const formik = useFormikContext<ResolverWizardFormProps>();
  const [showEditor, setShowEditor] = React.useState(false);
  const [editorValues, setEditorValues] = React.useState('');
  const [fileExtension, setFileExtension] =
    React.useState<SoloFormVisualEditorProps['mode']>('protobuf');
  const { setWarningMessage, warningMessage } = props;

  const updateExtension = (file: File) => {
    let extension = '';
    if (file.type) {
      extension = file.type;
    } else if (file.name) {
      extension = file.name.split('.').pop()!;
    }
    if (extension === 'proto') {
      return 'protobuf';
    } else if (extension === 'go') {
      return 'golang';
    } else if (extension === 'js') {
      return 'javascript';
    } else {
      return extension;
    }
  };

  return (
    <div data-testid='grpc-proto-check-section'>
      <Spacer px={6}>
        {/* TODO:  There's an edge case here where a user is adding multiple gRPC resolvers.
                   In this scenario, the proto page in the new resolver wizard should automatically
                   show the existing proto definition in the CR.
                   The user should have an option to update the proto bin if needed.
         */}
        <Spacer pb={2}>
          <WarningMessage message={warningMessage} />
        </Spacer>
        <SoloFormFileUpload
          data-testid='resolver-wizard-upload-proto'
          name='uploadProto'
          title='Upload gRPC pb.* file'
          fileType='*'
          buttonLabel='Upload ProtoFile'
          onRemoveFile={() => {
            setWarningMessage('');
            formik.setFieldValue('protoFile', '');
            setEditorValues('');
            setShowEditor(false);
          }}
          validateFile={file => {
            return new Promise((resolve, reject) => {
              if (file) {
                const reader = new FileReader();
                reader.onload = e => {
                  if (typeof e.target?.result === 'string') {
                    setFileExtension(updateExtension(file));
                    const proto = e.target.result;
                    try {
                      const isEncoded64 = isBase64(proto);
                      let newEditorValues = proto;
                      if (isEncoded64) {
                        newEditorValues = base64ToString(proto);
                        formik.setFieldValue('protoFile', proto);
                      } else {
                        newEditorValues = proto;
                        formik.setFieldValue(
                          'protoFile',
                          stringToBase64(proto)
                        );
                      }
                      setWarningMessage('');
                      setEditorValues(newEditorValues);
                      setShowEditor(true);
                      resolve({ isValid: true, errorMessage: '' });
                    } catch (err: any) {
                      setWarningMessage(err.message);
                      setEditorValues('');
                      setShowEditor(false);
                      reject({ isValid: false, errorMessage: err.message });
                    }
                  }
                };
                reader.readAsText(file);
              } else {
                resolve({ isValid: true, errorMessage: '' });
              }
            });
          }}
        />
      </Spacer>
      {showEditor && (
        <div className='mt-4'>
          <VisualEditor
            theme='chrome'
            name='grpcProtoCheck'
            style={{
              width: '100%',
              maxHeight: '46vh',
            }}
            focus={true}
            fontSize={16}
            showPrintMargin={false}
            showGutter={true}
            highlightActiveLine={true}
            defaultValue={editorValues}
            value={editorValues}
            readOnly={true}
            setOptions={{
              highlightGutterLine: true,
              showGutter: true,
              fontFamily: 'monospace',
              enableBasicAutocompletion: true,
              enableLiveAutocompletion: true,
              showLineNumbers: true,
              tabSize: 2,
            }}
            mode={fileExtension}
          />
        </div>
      )}
    </div>
  );
};
