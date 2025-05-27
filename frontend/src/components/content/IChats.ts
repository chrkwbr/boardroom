import {
  ApiErrorResult,
  ApiSuccessResult,
  del,
  get,
  post,
} from "../../fetch.ts";

export interface IChat {
  id: string | null;
  sender: string;
  image: string;
  message: string;
  version: number;
  date: Date;
}

export interface IPostChat {
  id: string | null;
  sender: string;
  room: string;
  message: string;
}

type EventType = "chat_created" | "chat_edited" | "chat_deleted";

export interface ChatEvent {
  id: string;
  event_type: EventType;
  sender: string;
  room: string;
  message: string;
  version: number;
  timestamp: number;
}

export const fetchChats: (roomId: string) => Promise<IChat[]> = async (
  roomId: string,
) => {
  interface IChatResponse {
    id: string;
    sender: string;
    message: string;
    version: number;
    date: number;
  }

  const apiResult: ApiSuccessResult<IChatResponse[]> | ApiErrorResult =
    await get<
      IChatResponse[]
    >(`chats/${roomId}/`);

  if (apiResult.ok && apiResult.data) {
    return Array.from(apiResult.data).map((it) => {
      return {
        id: it.id,
        sender: it.sender,
        image: "https://img.daisyui.com/images/profile/demo/1@94.webp",
        message: it.message,
        version: it.version,
        date: new Date(it.date * 1000),
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
  interface IChatHistoryResponse {
    id: string;
    sender: string;
    message: string;
    version: number;
    date: number;
  }

  const apiResult: ApiSuccessResult<IChatHistoryResponse[]> | ApiErrorResult =
    await get<
      IChatHistoryResponse[]
    >(`chats/${roomId}/${chatId}/history/`);

  if (apiResult.ok) {
    return Array.from(apiResult.data).map((it) => {
      return {
        id: it.id,
        sender: it.sender,
        image: "https://img.daisyui.com/images/profile/demo/1@94.webp",
        message: it.message,
        version: it.version,
        date: new Date(it.date * 1000),
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
