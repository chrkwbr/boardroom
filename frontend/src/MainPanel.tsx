import Header from "./header/Header.tsx";
import Content from "./content/Content.tsx";

const MainPanel = () => {
  return (
    <div className="flex flex-col w-full h-screen">
      <Header />
      <Content />
    </div>
  );
};

export default MainPanel;
