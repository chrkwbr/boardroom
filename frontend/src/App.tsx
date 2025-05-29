import "./App.css";
import Header from "./organisms/header/Header.tsx";
import { WebSocketProvider } from "./util/WebSocketProvider.tsx";
import Content from "./organisms/content/Content.tsx";

const App = () => {
  return (
    <WebSocketProvider>
      <div className="flex flex-col h-screen" data-theme="">
        <div className="flex-none">
          <Header />
        </div>
        <div className="flex flex-1 overflow-y-auto">
          <Content />
        </div>
      </div>
    </WebSocketProvider>
  );
};

export default App;
