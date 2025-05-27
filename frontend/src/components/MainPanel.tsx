import Sidebar from "./Sidebar.tsx";
import Content from "./content/Content.tsx";
import Header from "./header/Header.tsx";
import { Route, Routes } from "react-router-dom";

const MainPanel = () => {
  return (
    <div className="flex flex-col h-screen">
      <div className="flex-none">
        <Header />
      </div>
      <div className="flex flex-1 overflow-y-auto">
        <div className="flex sm:flex-col md:flex-row w-full">
          <Sidebar />
          <div className="px-1 flex-grow flex-shrink">
            <Routes>
              <Route
                path="/:roomId"
                element={<Content />}
              />
            </Routes>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MainPanel;
