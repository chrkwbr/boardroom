import MainPanel from "./MainPanel.tsx";

const Sidebar = () => {
  return (
    <div className="drawer lg:drawer-open">
      <input id="my-drawer" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content flex flex-col items-center justify-center">
        <MainPanel />
      </div>
      <div className="drawer-side bg-base-200 overflow-auto">
        <div className="p-3">
          <h1 className="mx-1 font-bold">
            <span className="text-3xl">Boardroom</span>
          </h1>
        </div>
        <label
          htmlFor="my-drawer"
          aria-label="close sidebar"
          className="drawer-overlay"
        >
        </label>
        <ul className="menu text-base-content w-80 p-4">
          <li>
            <a>Sidebar Item 1</a>
          </li>
          <li>
            <a>Sidebar Item 2</a>
          </li>
        </ul>
      </div>
    </div>
  );
};

export default Sidebar;
