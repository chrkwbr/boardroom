import { useState } from "react";

const ChatForm = (
  props: { onSend: (chat: string) => void; defaultText: string },
) => {
  const [chat, setChat] = useState<string>(props.defaultText);

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    handleSubmit();
  };

  const handleChatChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setChat(e.target.value);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleSubmit = () => {
    if (chat.trim() === "") {
      return;
    }
    props.onSend(chat);
    setChat("");
  };

  return (
    <div className="p-2">
      <form onSubmit={onSubmit}>
        <div className="flex flex-col border border-primary/50 rounded-lg shadow-lg">
          <div className="flex-1 justify-between items-center">
            <textarea
              className="textarea textarea-primary w-full"
              onChange={handleChatChange}
              onKeyDown={handleKeyDown}
              value={chat}
            >
            </textarea>
          </div>
          <div className="flex justify-end">
            <button className="btn btn-primary">Send</button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default ChatForm;
