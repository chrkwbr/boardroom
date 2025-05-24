const Sidebar = () => {
  return (
    <div className="bg-base-200 md:w-64 overflow-y-scroll sm:w-screen">
      <ul className="flex flex-col list-reset sm:hidden md:block">
        <li className="block">
          <a
            href="#"
            className="no-underline block h-full w-full px-8 py-4 hover:text-orange"
          >
            <i className="fa fa-tachometer mr-2" aria-hidden="true"></i>
            Dashboard
          </a>
        </li>
        <li className="flex justify-between">
          <a
            href="#"
            className="no-underline block h-full w-full px-8 py-4 hover:text-orange"
          >
            <i className="fa fa-user mr-2" aria-hidden="true"></i>
            Account
            <i className="fa fa-angle-right float-right" aria-hidden="true">
            </i>
          </a>
        </li>
        <li className="block">
          <a
            href="#"
            className="no-underline block h-full w-full px-8 py-4 hover:text-orange"
          >
            <i className="fa fa-envelope mr-2" aria-hidden="true"></i>
            MailBox
            <i className="fa fa-angle-down float-right" aria-hidden="true">
            </i>
          </a>
          <ul className="flex flex-col list-reset bg-orange-darkest">
            <li className="flex">
              <a
                href="#"
                className="no-underline block h-full w-full ml-4 hover:text-orange px-8 py-4"
              >
                <i className="fa fa-envelope-o mr-2" aria-hidden="true"></i>
                Inbox
              </a>
            </li>
            <li className="flex">
              <a
                href="#"
                className="no-underline block h-full w-full ml-4 hover:text-orange px-8 py-4"
              >
                <i className="fa fa-envelope-o mr-2" aria-hidden="true"></i>
                Categories
                <i
                  className="fa fa-angle-down float-right"
                  aria-hidden="true"
                >
                </i>
              </a>
            </li>
            <ul className="flex flex-col list-reset bg-orange-darkest">
              <li className="flex">
                <a
                  href="#"
                  className="no-underline block h-full w-full ml-8 hover:text-orange px-8 py-4"
                >
                  <i className="fa fa-envelope-o mr-2" aria-hidden="true"></i>
                  Social
                </a>
              </li>
              <li className="flex">
                <a
                  href="#"
                  className="no-underline block h-full w-full ml-8 hover:text-orange px-8 py-4"
                >
                  <i className="fa fa-envelope-o mr-2" aria-hidden="true"></i>
                  Notifications
                </a>
              </li>
            </ul>
          </ul>
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
