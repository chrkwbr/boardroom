import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { fetchRooms, IRoom } from "./IRooms.ts";
import { EventEmitter } from "../../util/EventEmitter.ts";
import { IChat } from "../chat/IChats.ts";

const Sidebar = () => {
  const [data, setData] = useState<IRoom[]>([]);

  useEffect(() => {
    const addChatListener = (event: { roomId: string; chat: IChat }) => {
      setData((prev) =>
        prev.map((room: IRoom) => {
          if (room.id === event.roomId) {
            return {
              ...room,
              unreadCount: (room.unreadCount || 0) + 1,
            };
          }
          return room;
        })
      );
    };
    EventEmitter.on("chat_created", addChatListener);

    (async () => {
      const rooms = await fetchRooms();
      if (!rooms) return;
      setData(rooms);
    })();

    return () => {
      EventEmitter.off("chat_created", addChatListener);
    };
  }, []);

  const handleRoomClick = (roomId: string) => {
    setData((prev: IRoom[]) =>
      prev.map((room: IRoom) => {
        if (room.id === roomId) {
          return { ...room, selected: true };
        }
        return { ...room, selected: false };
      })
    );
  };

  return (
    <div className="bg-base-200 md:w-52 overflow-y-scroll sm:w-screen">
      <ul className="menu rounded-box">
        <li>
          <details open>
            <summary>Parent</summary>
            <ul>
              {data.map((room) => (
                <li
                  key={room.id}
                  onClick={() => handleRoomClick(room.id)}
                  className={room.selected
                    ? "font-bold text-secondary bg-secondary-content"
                    : ""}
                >
                  <Link to={room.id}>
                    {room.name}{" "}
                    {room.unreadCount !== undefined && room.unreadCount > 0 && (
                      <div className="badge badge-xs badge-secondary">
                        {room.unreadCount}
                      </div>
                    )}
                  </Link>
                </li>
              ))}
            </ul>
          </details>
        </li>
        <li>
          <details open>
            <summary>Parent</summary>
            <ul>
              <li>
                <a>Submenu 1</a>
              </li>
              <li>
                <a>Submenu 2</a>
              </li>
              <li>
                <details open>
                  <summary>Parent</summary>
                  <ul>
                    <li>
                      <a>Submenu 1</a>
                    </li>
                    <li>
                      <a>Submenu 2</a>
                    </li>
                  </ul>
                </details>
              </li>
            </ul>
          </details>
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
