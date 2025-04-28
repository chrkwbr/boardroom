import {useCallback, useEffect, useState} from "react";
import ChatHistory from "./ChatHistory.tsx";
import ChatForm from "./ChatForm.tsx";
import {fetchChats, IChat} from "./IChats.ts";
import Header from "../Header.tsx";

const Content = () => {
  const [data, setData] = useState<IChat[]>([]);

  useEffect(() => {
    (async () => {
      const d: IChat[] = await fetchChats();
      setData(d);
    })();
  }, []);

  const handleSend = useCallback((chat: string) => {
    const newChat: IChat = {
      id: data.length,
      name: "You",
      image: "https://img.daisyui.com/images/profile/demo/1@94.webp",
      message: chat,
      date: new Date(),
    };
    setData((p: IChat[]) => [...p, newChat]);
  }, [data]);

  return (
    <div className="flex flex-col w-full h-screen">
      <Header />
      <div className="flex flex-col flex-1 bg-base-100">
        <div className="h-0 flex-1 overflow-y-auto" style={{ maxHeight: "calc(100vh - 10rem)"}}>
          <ChatHistory data={data} />
        </div>
        <div className="sticky bottom-0 left-0 right-0 bg-base-100">
          <ChatForm onSend={handleSend} />
        </div>
      </div>
    </div>
  );
};

export default Content;
