import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

type Room = {
  id: string;
  name: string;
  url: string;
};

const Sidebar = () => {
  const [data, setData] = useState<Room[]>([]);

  useEffect(() => {
    setData([
      {
        id: "1",
        name: "Room 1",
        url: "/room1",
      },
      {
        id: "2",
        name: "Room 2",
        url: "/room2",
      },
    ]);
  }, []);

  return (
    <div className="bg-base-200 md:w-52 overflow-y-scroll sm:w-screen">
      <ul className="menu rounded-box">
        <li>
          <details open>
            <summary>Parent</summary>
            <ul>
              {data.map((room) => (
                <li key={room.id}>
                  <Link to={room.url}>{room.name}</Link>
                </li>
              ))}
            </ul>
          </details>
        </li>
        <li>
          <details open>
            <summary>Parent</summary>
            <ul>
              <li>
                <a>Submenu 1</a>
              </li>
              <li>
                <a>Submenu 2</a>
              </li>
              <li>
                <details open>
                  <summary>Parent</summary>
                  <ul>
                    <li>
                      <a>Submenu 1</a>
                    </li>
                    <li>
                      <a>Submenu 2</a>
                    </li>
                  </ul>
                </details>
              </li>
            </ul>
          </details>
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
