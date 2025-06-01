import Sidebar from "../sidebar/Sidebar.tsx";
import { Route, Routes } from "react-router-dom";
import ChatRoom from "../chat/ChatRoom.tsx";
import { ChatEvent, IChat } from "../chat/IChats.ts";
import { useWebSocket } from "../../util/WebSocketProvider.tsx";
import { useEffect } from "react";
import { EventEmitter } from "../../util/EventEmitter.ts";

const Content = () => {
  const socket = useWebSocket();

  useEffect(() => {
    const chatHandler = (event: MessageEvent) => {
      try {
        const eventData = JSON.parse(event.data);
        const chatEvent = eventData satisfies ChatEvent;
        const chat: IChat = {
          id: chatEvent.id,
          sender: chatEvent.sender,
          image: "https://img.daisyui.com/images/profile/demo/1@94.webp",
          message: chatEvent.message,
          version: chatEvent.version,
          room: chatEvent.room,
          date: new Date(chatEvent.timestamp * 1000),
        };
        switch (eventData.event_type) {
          case "chat_created":
            EventEmitter.emit("chat_created", {
              roomId: chatEvent.room,
              chat: chat,
            });
            break;
          case "chat_edited":
            EventEmitter.emit("chat_edited", {
              roomId: chatEvent.room,
              chat: chat,
            });
            break;
          case "chat_deleted":
            EventEmitter.emit("chat_deleted", {
              roomId: chatEvent.room,
              chat: chat,
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
      <Sidebar />
      <div className="px-1 flex-grow flex-shrink">
        <Routes>
          <Route
            path="/:roomId"
            element={<ChatRoom />}
          />
        </Routes>
      </div>
    </div>
  );
};

export default Content;
