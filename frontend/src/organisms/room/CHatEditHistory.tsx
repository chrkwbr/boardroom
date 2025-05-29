import { useState } from "react";
import { fetchChatHistory, IChat } from "./IChats.ts";
import Chat from "./Chat.tsx";
import { useRoomId } from "./ChatRoom.tsx";

const CHatEditHistory = (props: { chatId: string }) => {
  const [data, setData] = useState<IChat[]>([]);
  const roomId = useRoomId();

  const handleShowHistory = () => {
    (async () => {
      const chatHistory: IChat[] = await fetchChatHistory(props.chatId, roomId);
      if (!chatHistory) return;
      setData(chatHistory);
    })();
  };

  return (
    <div className="drawer drawer-end">
      <input
        id={"my-drawer-" + props.chatId}
        type="checkbox"
        className="drawer-toggle"
      />
      <div className="px-1">
        <span
          className="text-xs text-secondary bg-neutral drawer-content"
          onClick={handleShowHistory}
        >
          <label
            htmlFor={"my-drawer-" + props.chatId}
            className="cursor-pointer"
          >
            edited
          </label>
        </span>
      </div>
      <div className="drawer-side" style={{ top: "4rem" }}>
        <label
          htmlFor={"my-drawer-" + props.chatId}
          aria-label="close sidebar"
          className="drawer-overlay"
        >
        </label>
        <ul className="menu min-h-full w-120 p-4 list rounded-box">
          {data.map((chat: IChat) => {
            return (
              <li key={chat.id + "-" + chat.version} className="list-row">
                <Chat chat={chat} />
              </li>
            );
          })}
        </ul>
      </div>
    </div>
  );
};

export default CHatEditHistory;
