import ChatHistory from "./ChatHistory.tsx";
import ChatForm from "./ChatForm.tsx";
import {fetchChats, IChat, postChat} from "./IChats.ts";
import {useCallback, useEffect, useRef, useState} from "react";

const Content = () => {
  const [data, setData] = useState<IChat[]>([]);
  const socketRef = useRef<WebSocket>(null);
  const [socketConnected, setSocketConnected] = useState(false);

  useEffect(() => {
    (async () => {
      const d: IChat[] = await fetchChats();
      setData(d);
    })();

    if (socketRef.current) {
      console.log("既存のWebSocket接続をクローズ中...");
      socketRef.current.close();
    }

    console.log("WebSocket接続を開始します...");
    const socket = new WebSocket("ws://localhost:8080/ws/chats/chnnelId1/");
    socketRef.current = socket;

    socket.addEventListener("open", () => {
      console.log("WebSocket接続確立");
      setSocketConnected(true);
    });

    socket.addEventListener("message", (event: MessageEvent) => {
      try {
        const eventData = JSON.parse(event.data);
        const newChat = eventData satisfies IChat;
        setData((prevData) => {
          if (newChat.id && prevData.some((chat) => chat.id === newChat.id)) {
            console.log("skip duplicated", newChat.id);
            return prevData;
          }
          return [...prevData, newChat];
        });
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

  const handleSend = useCallback((chat: string) => {
    (async () => {
      const newChat: IChat = {
        id: null,
        name: "You",
        image: "https://img.daisyui.com/images/profile/demo/1@94.webp",
        message: chat,
        date: new Date(),
      };
      await postChat(newChat);
    })();
  }, [data]);

  return (
    <div className="flex flex-col flex-1 bg-base-100">
      <div
        className="h-0 flex-1 overflow-y-auto"
        style={{ maxHeight: "calc(100vh - 10rem)" }}
      >
        <ChatHistory data={data} />
      </div>
      <div className="sticky bottom-0 left-0 right-0 bg-base-100">
        <ChatForm onSend={handleSend} />
      </div>
    </div>
  );
};

export default Content;
