import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar.jsx";
import { useAuth } from "../utils/AuthContext.jsx";
import ImageSlider from "../components/ProductImageSlider.jsx";
import Modal from "../components/Modal.jsx";

import PlusIcon from "../assets/plus.svg";
import DashIcon from "../assets/dash.svg";

import { useGetJson } from "../utils/customHook.js";

function ProductPage() {
  const { nm } = useParams();
  const [quantity, setQuantity] = useState(1);
  const { data: response, error: error } = useGetJson(`product/search/${nm}`);
  const { isAuth } = useAuth();
  const [first, setFirst] = useState(false);
  const [second, setSecond] = useState(false);

  const closeModalfirst = () => {
    setFirst(false);
  };

  const openModalfirst = () => {
    setFirst(true);
    setSecond(false);
  };

  const closeModalsecond = () => {
    setSecond(false);
  };

  const openModalsecond = () => {
    setSecond(false);
    setFirst(false);
  };

  const increment = () => {
    setQuantity((prev) => prev + 1);
  };

  const decrease = () => {
    quantity === 1 ? "" : setQuantity((prev) => prev - 1);
  };

  return (
    <>
      <Navbar />
      {response && response.data ? (
        <div className="max-w-[1400px] mx-auto h-[1000px]">
          <div className="max-w-[1000px] h-[650px]  mx-auto rounded-xl overflow-hidden mt-[50px] justify-center bg-slate-300">
            <div className="w-full h-full   justify-center">
              <ImageSlider itemID={response.data.id} />
            </div>
          </div>
          <div className="flex flex-col ml-[200px] bg-slate-300 rounded-xl  my-[10px] mr-[200px] p-2 font-semibold">
            <div className="flex text-[30px]">
              <h1 className="opacity-[0.90]">{response.data.name}</h1>
            </div>
            <div className="flex opacity-[0.75]">
              <p className="">{response.data.description}</p>
            </div>
            <div className="flex">
              <p className="">Rp.{response.data.price}</p>
            </div>
            <div className="flex h-[30px]">
              <div
                className={`border-2 border-rose-400 rounded-l-md w-[25px] ${
                  quantity === 1 ? "cursor-not-allowed" : "cursor-pointer"
                }`}
                onClick={decrease}
              >
                <img src={DashIcon} alt="not found images" className="" />
              </div>
              <div className="border-2 border-rose-400  w-[35px]">
                <p className="ml-2">{quantity}</p>
              </div>
              <div
                className={`border-2 border-rose-400 rounded-r-md w-[25px] cursor-pointer`}
                onClick={increment}
              >
                {/* <p className="ml-2">{quantity}</p> */}
                <img src={PlusIcon} alt="not found images" className="" />
              </div>
              <div className="flex">
                <p className="ml-[5px] border-2 border-slate-400 rounded-md">
                  Total Harga:{" "}
                  <span className="">Rp.{quantity * response.data.price}</span>{" "}
                </p>
                <button className=" border-2 border-sky-400 ml-[10px] rounded-md  hover:bg-slate-400">
                  Order Sekarang
                </button>
                <button className=" border-2 border-sky-400 ml-[10px] rounded-md hover:bg-slate-400">
                  Tambah Ke Keranjang
                </button>
              </div>
            </div>
          </div>
        </div>
      ) : (
        error && (
          <Modal
            title="Error Not Found"
            isOpen={true}
            content={<div className="mt-[5px]">cannot find product name</div>}
          />
        )
      )}
    </>
  );
}

export default ProductPage;
