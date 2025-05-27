const Dialog = (
  props: { id: string; text: string; title: string; deleteHandler: () => void },
) => {
  const handleDelete = () => {
    props.deleteHandler();
  };

  const dialogId = `my_modal_${props.id}`;

  return (
    <dialog id={dialogId} className="modal">
      <div className="modal-box">
        <h3 className="font-bold text-lg">{props.title}</h3>
        <p className="py-4">{props.text}</p>
        <div className="modal-action">
          <form method="dialog">
            <button className="btn">Cancel</button>
            <button className="btn btn-warning" onClick={handleDelete}>
              Delete
            </button>
          </form>
        </div>
      </div>
    </dialog>
  );
};

export default Dialog;
