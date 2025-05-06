import {ApiErrorResult, ApiSuccessResult, get, post} from "../../fetch.ts";

export interface IChat {
  id: string | null;
  sender: string;
  image: string;
  message: string;
  date: Date;
}

export interface IPostChat {
  id: string | null;
  sender: string;
  message: string;
}

export interface ChatEvent {
  ID: string | null;
  Sender: string;
  Room: string;
  Message: string;
  Timestamp: number;
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

export const postChat = async (chat: IPostChat) => {
  const apiResult: ApiSuccessResult<IPostChat> | ApiErrorResult = await post<
    IPostChat,
    IPostChat
  >("chats/channel/", chat);

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
};
