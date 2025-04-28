import {ApiErrorResult, ApiSuccessResult, get, post} from "../fetch.ts";

export interface IChat {
  id: string | null;
  name: string;
  image: string;
  message: string;
  date: Date;
}

export const fetchChats: () => Promise<IChat[]> = async () => {
  const apiResult: ApiSuccessResult<IChat[]> | ApiErrorResult = await get<
    IChat[]
  >("chats/channel/");

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    return [];
  }
};

export const fetchDiffChats: () => Promise<IChat[]> = async () => {
  const apiResult: ApiSuccessResult<IChat[]> | ApiErrorResult = await get<
    IChat[]
  >("chats/channel/diff/");

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    return [];
  }
};

export const postChat = async (chat: IChat) => {
  const apiResult: ApiSuccessResult<IChat> | ApiErrorResult = await post<
    IChat,
    IChat
  >("chats/channel/", chat);

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
};
