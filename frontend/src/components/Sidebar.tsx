import { NavLink } from "react-router-dom";
import { FiHome, FiSend, FiList, FiLayers, FiPieChart, FiBookOpen, FiUsers } from "react-icons/fi";

const Sidebar = () => {
  const navItems = [
    { name: "Dashboard", path: "/dashboard", icon: <FiHome className="nav-link-icon" /> },
    { name: "Send Money", path: "/send", icon: <FiSend className="nav-link-icon" /> },
    { name: "Transactions", path: "/transactions", icon: <FiList className="nav-link-icon" /> },
    { name: "Block Explorer", path: "/blocks", icon: <FiLayers className="nav-link-icon" /> },
    { name: "Reports", path: "/reports", icon: <FiPieChart className="nav-link-icon" /> },
    { name: "System Logs", path: "/logs", icon: <FiBookOpen className="nav-link-icon" /> },
    { name: "Beneficiaries", path: "/beneficiaries", icon: <FiUsers className="nav-link-icon" /> },
  ];

  return (
    <div className="sidebar-nav">
      {navItems.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={({ isActive }) =>
            `nav-link ${isActive ? "nav-link-active" : ""}`
          }
        >
          {item.icon}
          <span>{item.name}</span>
        </NavLink>
      ))}
    </div>
  );
};

export default Sidebar;
