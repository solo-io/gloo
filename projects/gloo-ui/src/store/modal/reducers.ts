import { MessageActionTypes, MessageAction } from './types';

export interface ModalState {
  message: string;
  error?: Error | null;
  showModal?: boolean;
}

const initialState: ModalState = {
  message: '',
  error: null,
  showModal: false
};

export function modalReducer(
  state = initialState,
  action: MessageActionTypes
): ModalState {
  switch (action.type) {
    case MessageAction.SUCCESS_MESSAGE:
      return {
        ...state,
        message: action.message,
        showModal: true
      };
    case MessageAction.HIDE_MODAL:
      return {
        ...state,
        message: '',
        showModal: false
      };
    case MessageAction.ERROR_MESSAGE:
      return {
        ...state,
        error: action.error,
        message: action.message
      };
    default:
      return {
        ...state,
        showModal: false
      };
  }
}
