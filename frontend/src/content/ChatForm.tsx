import {useState} from "react";

const ChatForm = (props: { onSend: (chat: string) => void }) => {
  const [chat, setChat] = useState<string>("");

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (chat.trim() === "") {
      return;
    }
    props.onSend(chat);
    setChat("");
  };

  const handleChatChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setChat(e.target.value);
  };

  return (
    <div className="w-full py-5">
      <form className="w-full" onSubmit={onSubmit}>
        <div className="flex justify-center">
          <div className="grid grid-cols-12 w-11/12">
            <div className="col-span-10">
              <textarea
                className="textarea textarea-secondary w-full"
                onChange={handleChatChange}
                value={chat}
              >
              </textarea>
            </div>
            <div className="col-span-2">
              <button className="btn btn-primary">Send</button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default ChatForm;
