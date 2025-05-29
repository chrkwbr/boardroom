import Sidebar, { SidebarChatHandlers } from "../sidebar/Sidebar.tsx";
import { Route, Routes } from "react-router-dom";
import ChatRoom, { ChatHandlers } from "../room/ChatRoom.tsx";
import { ChatEvent, IChat } from "../room/IChats.ts";
import { useWebSocket } from "../../util/WebSocketProvider.tsx";
import { useEffect, useRef } from "react";

const Content = () => {
  const socket = useWebSocket();
  const roomHandlers = useRef<
    Record<string, ChatHandlers>
  >({});
  const sidebarHandlers = useRef<SidebarChatHandlers>(null);

  const registerRoomHandlers = (roomId: string, handlers: ChatHandlers) => {
    roomHandlers.current[roomId] = handlers;
  };

  const registerSidebarHandlers = (handler: SidebarChatHandlers) => {
    sidebarHandlers.current = handler;
  };

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
          date: new Date(chatEvent.timestamp * 1000),
        };
        switch (eventData.event_type) {
          case "chat_created":
            roomHandlers.current[chatEvent.room]?.addChat(chat);
            sidebarHandlers.current?.addChat(chat);
            break;
          case "chat_edited":
            roomHandlers.current[chatEvent.room]?.editChat(chat);
            break;
          case "chat_deleted":
            roomHandlers.current[chatEvent.room]?.deleteChat(chat);
            break;
        }
      } catch (error) {
        console.error("Error parsing chat event data:", error);
      }
    };

    socket.addEventListener("message", chatHandler);
    return () => {
      console.log("Removing chat handler");
      socket.removeEventListener("message", chatHandler);
    };
  }, [socket]);

  return (
    <div className="flex sm:flex-col md:flex-row w-full">
      <Sidebar onRegister={registerSidebarHandlers} />
      <div className="px-1 flex-grow flex-shrink">
        <Routes>
          <Route
            path="/:roomId"
            element={<ChatRoom onRegister={registerRoomHandlers} />}
          />
        </Routes>
      </div>
    </div>
  );
};

export default Content;
