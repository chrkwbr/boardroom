import {useEffect} from "react";
import Chat from "./Chat.tsx";
import {IChat} from "./IChats.ts";

const ChatHistory = (props: { data: IChat[] }) => {
  useEffect(() => {
    props.data.forEach((chat) => {
      console.log(`Chat ID: ${chat.id},  Message: ${chat.message}`);
    });
  }, []);

  return (
    <ul className="list rounded-box shadow-md">
      {props.data.map((chat: IChat) => {
        return (
          <li key={chat.id} className="list-row">
            <Chat chat={chat} />
          </li>
        );
      })}
    </ul>
  );
};

export default ChatHistory;
