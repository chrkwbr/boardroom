import {useCallback, useEffect, useState} from "react";
import ChatHistory from "./ChatHistory.tsx";
import ChatForm from "./ChatForm.tsx";
import {fetchChats, IChat} from "./IChats.ts";

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
    <>
      <ChatHistory data={data} />
      <ChatForm onSend={handleSend} />
    </>
  );
};

export default Content;
