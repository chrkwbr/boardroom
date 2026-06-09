import {
  ApiErrorResult,
  ApiSuccessResult,
  del,
  get,
  post,
} from "../../util/fetch.ts";

export interface IChat {
  id: string | null;
  sender: ISender;
  roomId: string;
  message: string;
  version: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface ISender {
  id: string;
  name: string;
  icon: string;
}

export interface IPostChat {
  id: string | null;
  sender: string;
  room: string;
  message: string;
}

type EventType = "created" | "updated" | "deleted";

export interface ChatEvent {
  id: string;
  event_type: EventType;
  sender: string;
  room: string;
  message: string;
  version: number;
  timestamp: number;
}

interface IChatResponse {
  id: string | null;
  sender: ISender;
  roomId: string;
  message: string;
  version: number;
  createdAt: number
  updatedAt: number;
}

export const fetchChats: (roomId: string) => Promise<IChat[]> = async (
  roomId: string,
) => {
  const apiResult: ApiSuccessResult<IChatResponse[]> | ApiErrorResult =
    await get<
      IChatResponse[]
    >(`chats/${roomId}/`);

  if (apiResult.ok && apiResult.data) {
    return Array.from(apiResult.data).map((it) => {
      return {
        id: it.id,
        roomId: it.roomId,
        sender: it.sender,
        message: it.message,
        version: it.version,
        createdAt: new Date(it.createdAt),
        updatedAt: new Date(it.updatedAt),
      } as IChat;
    });
  } else {
    return [];
  }
};

export const fetchChatHistory = async (
  chatId: string,
  roomId: string,
): Promise<IChat[]> => {
  const apiResult: ApiSuccessResult<IChatResponse[]> | ApiErrorResult =
    await get<IChatResponse[]>(`chats/${roomId}/${chatId}/history/`);

  if (apiResult.ok) {
    return Array.from(apiResult.data).map((it) => {
      return {
        id: it.id,
        roomId: it.roomId,
        sender: it.sender,
        message: it.message,
        version: it.version,
        createdAt: new Date(it.createdAt),
        updatedAt: new Date(it.updatedAt),
      } as IChat;
    });
  } else {
    return [];
  }
};

export const postChat = async (chat: IPostChat) => {
  const apiResult: ApiSuccessResult<IPostChat> | ApiErrorResult = await post<
    IPostChat,
    IPostChat
  >(`chats/${chat.room}/`, chat);

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
  >(`chats/${chat.room}/${chat.id}`, chat);

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
  >(`chats/${chat.room}/${chat.id}`);

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
};
