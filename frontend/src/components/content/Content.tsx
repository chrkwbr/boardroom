import Chats from "./Chats.tsx";
import ChatForm from "./ChatForm.tsx";
import {ChatEvent, fetchChats, IChat, IPostChat, postChat} from "./IChats.ts";
import {useCallback, useEffect, useRef, useState} from "react";

const Content = () => {
  const [data, setData] = useState<IChat[]>([]);
  const dataRef = useRef<IChat[]>([]);
  const socketRef = useRef<WebSocket>(null);
  const [socketConnected, setSocketConnected] = useState(false);

  useEffect(() => {
    (async () => {
      const d: IChat[] = await fetchChats();
      if (!d) return;
      dataRef.current = d;
      setData(d);
    })();

    if (socketRef.current) {
      console.log("既存のWebSocket接続をクローズ中...");
      socketRef.current.close();
    }

    console.log("WebSocket接続を開始します...");
    const socket = new WebSocket("ws://localhost:8080/ws/chats/channel/");
    socketRef.current = socket;

    socket.addEventListener("open", () => {
      console.log("WebSocket接続確立");
      setSocketConnected(true);
    });

    socket.addEventListener("message", (event: MessageEvent) => {
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
        }
        switch (eventData.event_type) {
          case "chat_created":
            addChat(chat);
            break
          case "chat_edited":
            editChat(chat);
            break
          case "chat_deleted":
            deleteChat(chat);
            break
        }
      } catch (error) {
        console.error("メッセージの解析に失敗:", error);
      }
    });

    socket.addEventListener("error", (event) => {
      console.error("WebSocketエラー:", event);
    });

    socket.addEventListener("close", () => {
      console.log("closed WebSocket");
      setSocketConnected(false);
    });

    return () => {
      if (
        socketRef.current && socketRef.current.readyState === WebSocket.OPEN
      ) {
        socketRef.current.close();
      }
    };
  }, []);

  const addChat = (chat: IChat) => {
    if (chat.id && dataRef.current.some((c) => c.id === chat.id)) {
      console.log("skip duplicated", chat.id);
      return
    }
    const updatedData = [...dataRef.current, chat];
    dataRef.current = updatedData;
    setData(updatedData);
  }

  const editChat = (chat: IChat) => {
    const index = dataRef.current.findIndex((c) => c.id === chat.id);
    if (index === -1) {
      console.log("skip not found", chat.id);
      return
    }
    const updatedData = [...dataRef.current];
    updatedData[index] = chat;
    dataRef.current = updatedData;
    setData([...updatedData]);
  }

  const deleteChat = (chat: IChat) => {
    const index = dataRef.current.findIndex((c) => c.id === chat.id);
    if (index === -1) {
      console.log("skip not found", chat.id);
      return
    }
    const updatedData = [...dataRef.current];
    updatedData.splice(index, 1);
    dataRef.current = updatedData;
    setData([...updatedData]);
  }

  const handleSend = useCallback((chat: string) => {
    (async () => {
      const newChat: IPostChat = {
        id: null,
        sender: "You",
        message: chat,
      };
      await postChat(newChat);
    })();
  }, [data]);

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto">
        <Chats data={data} />
      </div>
      <div className="flex-none">
        <ChatForm onSend={handleSend} defaultText="" />
      </div>
    </div>
  );
};

export default Content;
