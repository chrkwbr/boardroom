import {ApiErrorResult, ApiSuccessResult, del, get, post} from "../../fetch.ts";

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

type EventType = "chat_created" | "chat_edited" | "chat_deleted";

export interface ChatEvent {
  id: string
  event_type: EventType
  sender: string
  room: string
  message: string
  timestamp: number
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

export const updateChat = async (chat: IPostChat) => {
  const apiResult: ApiSuccessResult<IPostChat> | ApiErrorResult = await post<
    IPostChat,
    IPostChat
  >(`chats/channel/${chat.id}`, chat);

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
};

export const deleteChat = async (chat: IPostChat) => {
  const apiResult: ApiSuccessResult<string> | ApiErrorResult = await del<
    string
  >(`chats/channel/${chat.id}`);

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
}