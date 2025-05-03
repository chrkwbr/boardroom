import Chat from "./Chat.tsx";
import {IChat} from "./IChats.ts";

const ChatHistory = (props: { data: IChat[] }) => {
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
