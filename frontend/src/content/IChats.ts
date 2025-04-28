import {ApiErrorResult, ApiSuccessResult, get, post} from "../fetch.ts";

export interface IChat {
  id: number;
  name: string;
  image: string;
  message: string;
  date: Date;
}

export const fetchChats: () => Promise<IChat[]> = async () => {
  const apiResult: ApiSuccessResult<IChat[]> | ApiErrorResult = await get<
    IChat[]
  >("channel/");

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
  >("channel/", chat);

  if (apiResult.ok) {
    console.log(`posted ${JSON.stringify(apiResult.data)} chat`);
    return apiResult.data;
  } else {
    alert(`Error: ${apiResult.data.message}`);
    return null;
  }
};
