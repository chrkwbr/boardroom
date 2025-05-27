import { ApiErrorResult, ApiSuccessResult, get } from "../../fetch.ts";

export interface INotification {
  id: number;
  message: string;
  date: Date;
  read: boolean;
}

export const getNotifications: () => Promise<INotification[]> = async () => {
  const apiResult: ApiSuccessResult<INotification[]> | ApiErrorResult =
    await get<
      INotification[]
    >("notifications/");

  if (apiResult.ok) {
    return apiResult.data;
  } else {
    return [];
  }
};
