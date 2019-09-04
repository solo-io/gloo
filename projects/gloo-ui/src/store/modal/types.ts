export enum MessageAction {
  SUCCESS_MESSAGE = 'SUCCESS_MESSAGE',
  ERROR_MESSAGE = 'ERROR_MESSAGE',
  SHOW_MODAL = 'SHOW_MODAL',
  HIDE_MODAL = 'HIDE MODAL'
}

export interface SuccessMessageAction {
  type: typeof MessageAction.SUCCESS_MESSAGE;
  message: string;
}

export interface ShowModalAction {
  type: MessageAction.SHOW_MODAL;
}

export interface HideModalAction {
  type: MessageAction.HIDE_MODAL;
}
export interface ErrorMessageAction {
  type: typeof MessageAction.ERROR_MESSAGE;
  error: Error | null;
  message: string;
}

export type MessageActionTypes =
  | SuccessMessageAction
  | ErrorMessageAction
  | HideModalAction;
