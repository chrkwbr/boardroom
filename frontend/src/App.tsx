import "./App.css";
import Header from "./components/header/Header.tsx";
import Sidebar from "./components/room/Sidebar.tsx";
import Content from "./components/content/Content.tsx";
import {Route, Routes} from "react-router-dom";

const App = () => {
  return (
    <div className="flex flex-col h-screen" data-theme="">
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

export default App;
