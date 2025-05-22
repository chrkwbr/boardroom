// deno-lint-ignore-file ban-ts-comment
import {useCallback, useState} from "react";
import {deleteChat, IChat, IPostChat, updateChat} from "./IChats.ts";
import ChatForm from "./ChatForm.tsx";
import Dialog from "../modal/dialog.tsx";

const Chat = (props: { chat: IChat }) => {
  const [editing, setEditing] = useState(false);

  const startEdit = () => {
    setEditing(true);
  };

  const cancelEdit = () => {
    setEditing(false);
  };

  const handleDelete = () => {
    const elementById = document.getElementById(`my_modal_${props.chat.id}`);
    if (elementById) {
      // @ts-ignore
      elementById.showModal();
    }
  };

  const delChat = () => {
    (async () => {
      const toDelete: IPostChat = {
        id: props.chat.id,
        sender: "You",
        message: "",
      };
      await deleteChat(toDelete);
    })();
  };

  const handleEdit = useCallback((chat: string) => {
    (async () => {
      const newChat: IPostChat = {
        id: props.chat.id,
        sender: "You",
        message: chat,
      };
      await updateChat(newChat);
    })();
    cancelEdit();
  }, []);

  return (
    <>
      <div>
        <img className="size-10 rounded-box" src={props.chat.image} />
      </div>
      <div>
        <div>{props.chat.sender}</div>
        <div className="text-xs uppercase font-semibold opacity-60">
          {props.chat.date.toString()}
        </div>
      </div>
      {!editing &&
        (
          <p className="list-col-wrap text-xs">
            {props.chat.message}
          </p>
        )}
      <div className="list-col-wrap text-xs">
        {editing &&
          (
            <div className="w-full">
              <ChatForm onSend={handleEdit} defaultText={props.chat.message} />
            </div>
          )}
      </div>
      <button className="btn btn-square btn-ghost">
        <svg
          className="size-[1.2em]"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <g
            strokeLinejoin="round"
            strokeLinecap="round"
            strokeWidth="2"
            fill="none"
            stroke="currentColor"
          >
            <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z">
            </path>
          </g>
        </svg>
      </button>
      <div className="flex-none">
        <div className="dropdown dropdown-end">
          <div tabIndex={0} role="button" className="btn btn-ghost btn-circle">
            <div className="indicator">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                className="inline-block h-5 w-5 stroke-current"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z"
                >
                </path>
              </svg>
            </div>
          </div>
          <ul
            tabIndex={0}
            className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
          >
            {editing
              ? (
                <>
                  <li onClick={cancelEdit}>
                    <button type="button" className="btn btn-sm btn-outline btn-secondary">Cancel Edit</button>
                  </li>
                </>
              )
              : (
                <>
                  <li onClick={startEdit}>
                    <button type="button" className="btn btn-sm btn-outline btn-secondary">Edit</button>
                  </li>
                </>
              )}
            <li onClick={handleDelete}>
              <button type="button" className="btn btn-sm btn-outline btn-warning">Delete</button>
            </li>
          </ul>
        </div>
      </div>
      <Dialog
        deleteHandler={delChat}
        id={props.chat.id!}
        text={`Are you sure you want to delete message?`}
        title={`Delete message`}
      />
    </>
  );
};

export default Chat;
