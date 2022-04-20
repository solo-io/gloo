import ConfirmationModal from 'Components/Common/ConfirmationModal';
import React, { useState, useCallback, createContext, useContext } from 'react';

//
// TYPES
//
type ConfirmOptions = {
  confirmPrompt: string;
  confirmButtonText: string;
  isNegative?: boolean;
};

interface IConfirmModalContext {
  confirm(options: ConfirmOptions): Promise<any>;
}

//
// CONTEXT
//
const ConfirmModalContext = createContext({} as IConfirmModalContext);

//
// PROVIDER
//
export const ConfirmModalProvider: React.FC = props => {
  const [confirmOptions, setConfirmOptions] = useState<ConfirmOptions | null>(
    null
  );
  const [resolveReject, setResolveReject] = useState<any[]>([]);
  const [resolve, reject] = resolveReject;

  const confirm = useCallback((options: ConfirmOptions) => {
    return new Promise((resolve, reject) => {
      setConfirmOptions(options);
      setResolveReject([resolve, reject]);
    });
  }, []);

  const handleClose = useCallback(() => {
    setResolveReject([]);
  }, []);

  const handleCancel = useCallback(() => {
    if (reject) {
      reject();
      handleClose();
    }
  }, [reject, handleClose]);

  const handleConfirm = useCallback(() => {
    if (resolve) {
      resolve();
      handleClose();
    }
  }, [resolve, handleClose]);

  return (
    <ConfirmModalContext.Provider value={{ confirm }}>
      <ConfirmationModal
        visible={resolveReject.length === 2}
        goForIt={handleConfirm}
        cancel={handleCancel}
        confirmTestId='confirm-modal-button'
        confirmPrompt={confirmOptions?.confirmPrompt}
        confirmButtonText={confirmOptions?.confirmButtonText}
        isNegative={confirmOptions?.isNegative}
      />
      {props.children}
    </ConfirmModalContext.Provider>
  );
};

//
// HOOK
//
export const useConfirm = () => {
  const { confirm } = useContext(ConfirmModalContext);
  return confirm;
};
