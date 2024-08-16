import { useState } from "react";
import { useAuth } from "../utils/AuthContext.jsx";
import { getJson } from "../utils/instance.js";
import Modal from "./Modal.jsx";
import { Link, useLocation } from "react-router-dom";
import IoSearch from "../assets/search-outline.svg";

// eslint-disable-next-line react/prop-types
export default function Navbar() {
  const { isAuth } = useAuth();
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState("");

  const onChange = (e) => {
    setInputValue(e.target.value);
  };

  const openModal = () => {
    setIsOpen(true);
  };

  const closeModal = () => {
    setIsOpen(false);
    setInputValue("");
  };

  const logOutButton = async () => {
    try {
      const response = await getJson("signout");
      if (response.status === 200) {
        window.location.reload();
      } else {
        throw new Error("error response not 200:", response);
      }
    } catch (error) {
      alert(`ada error:${error}`);
    }
  };

  const getPathName = (path) => {
    return location.pathname == path && "text-slate-900 bg-slate-200";
  };

  const AuthValue = (
    <>
      <Link
        to="/signin"
        className={`font-bold px-3 py-2 text-slate-700 rounded-lg hover:bg-slate-100 hover:text-slate-900 ${getPathName(
          "/signin"
        )}`}
      >
        Signin
      </Link>
      <Link
        to="/signup"
        className={`font-bold px-3 py-2 text-slate-700 rounded-lg  hover:bg-slate-100 hover:text-slate-900 ${getPathName(
          "/signup"
        )}`}
      >
        Signup
      </Link>
    </>
  );
  return (
    <div>
      <nav className="flex justify-center space-x-4 mt-4 m-2 ">
        <Link
          to="/"
          className={`font-bold px-3 py-2 text-slate-700 rounded-lg  hover:bg-slate-100 hover:text-slate-900 ${getPathName(
            "/"
          )}`}
        >
          Home
        </Link>

        {isAuth ? (
          <>
            <button
              type="submit"
              onClick={logOutButton}
              className="font-bold px-3 py-2 text-slate-700 rounded-lg  hover:bg-slate-100 hover:text-slate-900  "
            >
              Signout
            </button>
            <Link
              to="/profile"
              className={`font-bold px-3 py-2 text-slate-700 rounded-lg  hover:bg-slate-100 hover:text-slate-900 ${getPathName(
                "/profile"
              )}`}
            >
              Profile
            </Link>
          </>
        ) : (
          AuthValue
        )}
        <div
          className="mt-1 hover:bg-slate-100 cursor-pointer"
          onClick={openModal}
        >
          <img src={IoSearch} alt="not found images" className="size-[35px]" />
        </div>
      </nav>
      <Modal
        isOpen={isOpen}
        title="Cari Produk"
        content={
          <div className="flex w-full mt-[15px] ">
            <input
              type="text"
              value={inputValue}
              onChange={onChange}
              className="border-2 border-slate-500 mr-[3px] rounded-l-lg focus:outline-none text-center"
            />
            <Link to={`/p/${inputValue}`} className="">
              <img
                src={IoSearch}
                alt="not found images"
                className="size-[35px] border-2 border-slate-500 rounded-r-lg cursor-pointer hover:bg-slate-400"
              />
            </Link>
          </div>
        }
        onCloseButton={
          <>
            <div className="border-2 border-slate-300 rounded-t-lg w-[60px] text-center hover:bg-green-300 mt-8">
              <button className="" onClick={closeModal}>
                Close
              </button>
            </div>
          </>
        }
      />
    </div>
  );
}
