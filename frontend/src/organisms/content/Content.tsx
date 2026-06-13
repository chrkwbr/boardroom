import Sidebar from "../sidebar/Sidebar.tsx";
import {Route, Routes} from "react-router-dom";
import ChatRoom from "../chat/ChatRoom.tsx";
import {IChat, IChatResponse} from "../chat/IChats.ts";
import {useWebSocket} from "../../util/WebSocketProvider.tsx";
import {useEffect} from "react";
import {EventEmitter} from "../../util/EventEmitter.ts";

interface IChatWsEvent {
  event_type: string;
  room_id: string;
  chat_id: string;
  chat: IChatResponse | null;
}

const Content = () => {
  const socket = useWebSocket();

  useEffect(() => {
    if (!socket) return;

    const chatHandler = (event: MessageEvent) => {
      try {
        const eventData = JSON.parse(event.data);
        const chatEvent = eventData satisfies IChatWsEvent;
        switch (chatEvent.event_type) {
          case "created":
            let createdChat = {
              id: chatEvent.chat.id,
              roomId: chatEvent.chat.roomId,
              sender: chatEvent.chat.sender,
              message: chatEvent.chat.message,
              version: chatEvent.chat.version,
              createdAt: new Date(chatEvent.chat.createdAt),
              updatedAt: new Date(chatEvent.chat.updatedAt),
            } as IChat
            EventEmitter.emit("chat_created", {
              roomId: chatEvent.chat.roomId,
              chat:  createdChat,
            });
            break;
          case "updated":
            const updatedChat: IChat = {
              id: chatEvent.chat.id,
              roomId: chatEvent.chat.roomId,
              sender: chatEvent.chat.sender,
              message: chatEvent.chat.message,
              version: chatEvent.chat.version,
              createdAt: new Date(chatEvent.chat.createdAt),
              updatedAt: new Date(chatEvent.chat.updatedAt),
            } as IChat
            EventEmitter.emit("chat_edited", {
              roomId: chatEvent.chat.roomId,
              chat: updatedChat,
            });
            break;
          case "deleted":
            EventEmitter.emit("chat_deleted", {
              roomId: chatEvent.room_id,
              chatId: chatEvent.chat_id
            });
            break;
        }
      } catch (error) {
        console.error("Error parsing chat event data:", error);
      }
    };

    socket.addEventListener("message", chatHandler);
    return () => {
      socket.removeEventListener("message", chatHandler);
    };
  }, [socket]);

  return (
    <div className="flex sm:flex-col md:flex-row w-full">
      <Sidebar/>
      <div className="px-1 flex-grow flex-shrink">
        <Routes>
          <Route
            path="/:roomId"
            element={<ChatRoom/>}
          />
        </Routes>
      </div>
    </div>
  );
};

export default Content;
