import { getNotifications, INotification } from "./INotification.ts";
import { useEffect, useState } from "react";

const Notifications = () => {
  const [notifications, setNotifications] = useState<INotification[]>([]);

  useEffect(() => {
    polling();
  }, []);

  const polling = () => {
    // setInterval(async () => {
    //   const newOne = await getNotifications();
    //   setNotifications((prev) => [...prev, ...newOne]);
    // }, 5000);
  };

  return (
    <div className="dropdown dropdown-end">
      <div tabIndex={0} role="button" className="btn btn-ghost btn-circle">
        <div className="indicator">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
            />
          </svg>
          {notifications.length > 0 &&
            (
              <span className="badge badge-xs badge-primary indicator-item">
              </span>
            )}
        </div>
      </div>
      <div
        tabIndex={0}
        className="card card-compact dropdown-content bg-secondary z-1 mt-3 w-96 shadow"
      >
        {notifications.length > 0 &&
          (
            <ul className="menu menu-sm rounded-box rounded-t-none p-2">
              {notifications.map((notification) => (
                <li key={notification.id}>
                  <a className="justify-between">
                    {notification.message}
                  </a>
                </li>
              ))}
            </ul>
          )}
      </div>
    </div>
  );
};

export default Notifications;
