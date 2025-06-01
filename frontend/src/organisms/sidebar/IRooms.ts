export interface IRoom {
  id: string;
  name: string;
  unreadCount?: number; // Optional property to track unread messages
  selected?: boolean; // Optional property to indicate if the room is selected
}

export const fetchRooms = async (): Promise<IRoom[]> => {
  interface IRoomResponse {
    id: string;
    name: string;
    url: string;
  }

  return [{
    id: "room1",
    name: "Room 1",
    selected: true,
  }, {
    id: "room2",
    name: "Room 2",
  }, {
    id: "room3",
    name: "Room 3",
  }];

  // const apiResult: ApiSuccessResult<IRoomResponse[]> | ApiErrorResult =
  //   await get<IRoomResponse[]>("rooms/");
  //
  // if (apiResult.ok && apiResult.data) {
  //   return Array.from(apiResult.data).map((it) => {
  //     return {
  //       id: it.id,
  //       name: it.name,
  //       url: it.url,
  //     } as IRoom;
  //   });
  // } else {
  //   return [];
  // }
};
