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
    <div className="py-2">
      <form onSubmit={onSubmit}>
        <div className="flex justify-center">
          <div className="grid grid-cols-12 w-full">
            <div className="col-span-10 px-1">
              <textarea
                className="textarea textarea-primary w-full"
                onChange={handleChatChange}
                onKeyDown={handleKeyDown}
                value={chat}
              >
              </textarea>
            </div>
            <div className="col-span-2 px-1">
              <button className="btn btn-primary">Send</button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default ChatForm;
