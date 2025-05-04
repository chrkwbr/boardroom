import Chat from "./Chat.tsx";
import {IChat} from "./IChats.ts";
import {useLayoutEffect, useRef} from "react";

const ChatHistory = (props: { data: IChat[] }) => {
  const endOfMessages = useRef<HTMLDivElement>(null);
  useLayoutEffect(() => {
    if (endOfMessages.current) {
      endOfMessages.current.scrollIntoView({
        behavior: "smooth",
        block: "end",
      });
    }
  });

  return (
    <ul className="list rounded-box shadow-md">
      {props.data.map((chat: IChat) => {
        return (
          <li key={chat.id} className="list-row">
            <Chat chat={chat} />
          </li>
        );
      })}
      <div ref={endOfMessages}></div>
    </ul>
  );
};

export default ChatHistory;
