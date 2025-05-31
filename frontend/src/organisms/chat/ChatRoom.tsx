import ChatTimeline from "./ChatTimeline.tsx";
import ChatForm from "./ChatForm.tsx";
import { fetchChats, IChat, IPostChat, postChat } from "./IChats.ts";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { useParams } from "react-router-dom";
import { EventEmitter } from "../../util/EventEmitter.ts";

export const RoomContext = createContext<string | null>(null);

export const useRoomId = () => {
  const roomId = useContext(RoomContext);
  if (!roomId) {
    throw new Error("RoomContext is not provided");
  }
  return roomId;
};

const ChatRoom = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const [data, setData] = useState<IChat[]>([]);
  const dataRef = useRef<IChat[]>([]);

  useEffect(() => {
    if (!roomId) return;
    const addChatListener = (event: { roomId: string; chat: IChat }) => {
      if (event.roomId !== roomId) return;
      addChat(event.chat);
    };
    EventEmitter.on("chat_created", addChatListener);

    const editChatListener = (event: { roomId: string; chat: IChat }) => {
      if (event.roomId !== roomId) return;
      editChat(event.chat);
    };
    EventEmitter.on("chat_edited", editChatListener);

    const deleteChatListener = (event: { roomId: string; chat: IChat }) => {
      if (event.roomId !== roomId) return;
      deleteChat(event.chat);
    };
    EventEmitter.on("chat_deleted", deleteChatListener);

    (async () => {
      const d: IChat[] = await fetchChats(roomId);
      if (!d) return;
      dataRef.current = d;
      setData(d);
    })();

    return () => {
      EventEmitter.off("chat_created", addChatListener);
      EventEmitter.off("chat_edited", editChatListener);
      EventEmitter.off("chat_deleted", deleteChatListener);
    };
  }, [roomId]);

  const addChat = (chat: IChat) => {
    if (chat.id && dataRef.current.some((c) => c.id === chat.id)) {
      console.log("skip duplicated", chat.id);
      return;
    }
    const updatedData = [...dataRef.current, chat];
    dataRef.current = updatedData;
    setData(updatedData);
  };

  const editChat = (chat: IChat) => {
    const index = dataRef.current.findIndex((c) => c.id === chat.id);
    if (index === -1) {
      console.log("skip not found", chat.id);
      return;
    }
    const updatedData = [...dataRef.current];
    updatedData[index] = chat;
    dataRef.current = updatedData;
    setData([...updatedData]);
  };

  const deleteChat = (chat: IChat) => {
    const index = dataRef.current.findIndex((c) => c.id === chat.id);
    if (index === -1) {
      console.log("skip not found", chat.id);
      return;
    }
    const updatedData = [...dataRef.current];
    updatedData.splice(index, 1);
    dataRef.current = updatedData;
    setData([...updatedData]);
  };

  const handleSend = useCallback((chat: string) => {
    (async () => {
      const newChat: IPostChat = {
        id: null,
        sender: "You",
        room: roomId!,
        message: chat,
      };
      await postChat(newChat);
    })();
  }, [data]);

  return (
    roomId
      ? (
        <RoomContext.Provider value={roomId}>
          <div className="flex flex-col h-full">
            <div className="flex-1 overflow-y-auto">
              <ChatTimeline data={data} />
            </div>
            <div className="flex-none">
              <ChatForm onSend={handleSend} defaultText="" />
            </div>
          </div>
        </RoomContext.Provider>
      )
      : <></>
  );
};

export default ChatRoom;
