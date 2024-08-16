import { useState } from "react";
import { useParams } from "react-router-dom";
import Navbar from "../components/Navbar.jsx";
import { useAuth } from "../utils/AuthContext.jsx";
import ImageSlider from "../components/ProductImageSlider.jsx";
import Modal from "../components/Modal.jsx";
import ModalLoginErrorMsg from "../components/ModalLoginErrorMsg.jsx";
import PlusIcon from "../assets/plus.svg";
import DashIcon from "../assets/dash.svg";

import { useGetJson, usePostJson } from "../utils/customHook.js";

function ProductPage() {
  const { nm } = useParams();
  const [quantity, setQuantity] = useState(1);
  const { data: response, error: error } = useGetJson(`product/search/${nm}`);
  const {
    data: orderResponse,
    error: orderError,
    call: orderCall,
    // resetResponse: orderReset,
  } = usePostJson("order/create");
  const {
    data: cartResponse,
    error: cartError,
    call: cartCall,
    // resetResponse: cartReset,
  } = usePostJson("cart/create");

  const { isAuth } = useAuth();
  const [orderModal, setOrderModal] = useState(false);
  const [cartModal, setCartModal] = useState(false);

  const closeModal = () => {
    setOrderModal(false);
    setCartModal(false);
  };

  const openOrderModal = () => {
    setOrderModal(true);
    setCartModal(false);
  };

  const openCartModal = () => {
    setCartModal(true);
    setOrderModal(false);
  };

  const increment = () => {
    quantity == response.data.stock
      ? setQuantity(response.data.stock)
      : quantity <= response.data.stock && setQuantity((prev) => prev + 1);
  };

  const decrease = () => {
    quantity === 1 ? "" : setQuantity((prev) => prev - 1);
  };

  const orderSubmit = async () => {
    const form = {
      product_id: response && response.data.id,
      total: quantity,
    };
    await orderCall({ form: form });
  };

  const cartSubmit = async () => {
    const form = {
      product_id: response && response.data.id,
      status: "wait",
      total: quantity,
    };
    await cartCall({ form: form });
  };

  function activeButton(id) {
    const text = document.getElementById(id);
    if (text.classList.contains("hidden")) {
      text.classList.remove("hidden");
    } else {
      text.classList.add("hidden");
    }
  }

  return (
    <>
      <Navbar />
      {response && response.data ? (
        <div className="max-w-[1400px] mx-auto h-[1000px]">
          <div className="max-w-[1000px] h-[650px]  mx-auto rounded-xl overflow-hidden mt-[50px] justify-center">
            <div className="w-full h-full   justify-center">
              <ImageSlider itemID={response.data.id} />
            </div>
          </div>
          <div className="flex flex-col ml-[200px]  rounded-xl  my-[10px] mr-[200px] p-2 font-semibold border-[1px] border-slate-400">
            <div className="text-[30px] underline underline-offset-8 mb-2">
              <h1 className="opacity-[0.90] ">{response.data.name}</h1>
            </div>
            {/* batas */}

            {/* <div className="flex border-y-[2px] border-slate-400 my-4 p-2">
              <div className="flex">
                  <button className="peer m-2 hover:text-sky-500 hover:underline hover:underline-offset-[22px] focus:text-sky-500 focus:underline focus:underline-offset-[22px] focus:decoration-blue-500 focus:decoration-[2px]">
                    Deskripsi
                  </button>

                <button className="m-2  hover:text-sky-500 hover:underline hover:underline-offset-[22px] focus:text-sky-500 focus:underline focus:underline-offset-[22px] focus:decoration-blue-800 focus:decoration-[2px]">
                  Diskusi
                </button>
              </div>
            </div>
            <div className="h-[50px] font-normal p-2  mb-[30px]">
              <p className="hidden peer-hover:block ">Deskripsi produk</p>
              <p className="hidden">Diskusi Produk</p>
            </div> */}

            {/* start */}
            <div className="border-y-[2px] border-slate-400 my-6">
              <div className="flex ">
                <button
                  onClick={() => activeButton("text1")}
                  className="peer m-2 hover:text-sky-500 hover:underline hover:underline-offset-[15px] active:text-sky-500 active:underline active:underline-offset-[14px] active:decoration-blue-500 focus:decoration-[2px]"
                >
                  Deskripsi
                </button>
                {/* <button className=" m-2 hover:text-sky-500 hover:underline hover:underline-offset-[15px] focus:text-sky-500 focus:underline focus:underline-offset-[14px] focus:decoration-blue-500 focus:decoration-[2px]">
                  Diskusi
                </button> */}
              </div>
            </div>

            <div className="h-[50px] font-normal p-2  mb-[30px] ">
              <p id="text1" className={`hidden`}>
                {response.data.description}
              </p>
            </div>

            {/* end */}

            {/* batas  */}

            <div className="flex mb-4">
              <p className="">Rp.{response.data.price}</p>
            </div>
            <div className="flex h-[30px]">
              <div
                className={`border-2 border-rose-400 rounded-l-md w-[25px] ${
                  quantity === 1 ? "cursor-not-allow3ed" : "cursor-pointer"
                }`}
                onClick={decrease}
              >
                <img src={DashIcon} alt="not found images" className="" />
              </div>
              <div className="border-2 border-rose-400  w-[50px]">
                <p className="mx-[10px]">{quantity}</p>
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
                {isAuth ? (
                  <Modal
                    title="Confirm Order"
                    isOpen={orderModal}
                    content={
                      <>
                        <div
                          className={`p-[10px] border-[1px] my-[15px] rounded-md ${
                            orderError ? "border-red-800" : "border-slate-400"
                          }`}
                        >
                          <p className="font-normal">
                            Pengiriman akan dikirim ke alamat mu
                          </p>
                          <h1 className="flex font-bold">
                            Jumlah order:{" "}
                            <p className="text-red-700 ml-[5px]">{quantity}</p>
                          </h1>
                          <h2 className="flex font-bold">
                            Harga yang harus dibayar:
                            <p className="text-green-700 ml-[5px]">
                              Rp.{quantity * response.data.price}
                            </p>
                          </h2>
                        </div>
                        <span
                          className={`text-sm m-2 font-semibold ${
                            orderError
                              ? "text-pink-800"
                              : orderResponse && "text-green-500"
                          } underline`}
                        >
                          {orderError
                            ? orderError.message
                            : orderResponse
                            ? !orderError && orderResponse.message
                            : ""}
                        </span>
                      </>
                    }
                    onCloseButton={
                      <div>
                        <button
                          onClick={closeModal}
                          className="transition hover:scale-[1.05] mt-12"
                        >
                          Close
                        </button>
                        <button
                          onClick={orderSubmit}
                          className="transition hover:scale-[1.05] mt-12 ml-[220px]"
                        >
                          Pesan
                        </button>
                      </div>
                    }
                  />
                ) : (
                  orderModal && (
                    <ModalLoginErrorMsg
                      closeButton={
                        <button
                          onClick={closeModal}
                          className="transition hover:scale-[1.05] "
                        >
                          Close
                        </button>
                      }
                    />
                  )
                )}

                <button
                  onClick={openOrderModal}
                  className=" border-2 border-sky-400 ml-[10px] rounded-md  hover:bg-slate-400"
                >
                  Order Sekarang
                </button>
                {isAuth ? (
                  <Modal
                    title="Keranjang"
                    isOpen={cartModal}
                    content={
                      <>
                        <div
                          className={`p-2 mt-[15px] mb-2 flex items-center justify-center container mx-auto border-[1px] ${
                            cartError ? "border-red-800" : "border-slate-300"
                          } rounded-md`}
                        >
                          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-4">
                            <div className="">
                              <ImageSlider itemID={response.data.id} />
                            </div>
                            <p className="font-normal">{response.data.name}</p>
                            <p className="font-semibold">{quantity}x</p>
                            <p className="font-bold">
                              Rp.{quantity * response.data.price}
                            </p>
                          </div>
                        </div>
                        <span
                          className={`text-sm font-semibold ${
                            cartError
                              ? "text-pink-800"
                              : cartResponse && "text-green-500"
                          } underline`}
                        >
                          {cartError
                            ? cartError.message
                            : cartResponse
                            ? !cartError && cartResponse.message
                            : ""}
                        </span>
                      </>
                    }
                    onCloseButton={
                      <div>
                        <button
                          onClick={closeModal}
                          className="transition hover:scale-[1.05] mt-12"
                        >
                          Close
                        </button>
                        <button
                          onClick={cartSubmit}
                          className="transition hover:scale-[1.05] mt-12 ml-[220px]"
                        >
                          Tambah
                        </button>
                      </div>
                    }
                  />
                ) : (
                  cartModal && (
                    <ModalLoginErrorMsg
                      closeButton={
                        <button
                          onClick={closeModal}
                          className="transition hover:scale-[1.05] "
                        >
                          Close
                        </button>
                      }
                    />
                  )
                )}
                <button
                  onClick={openCartModal}
                  className=" border-2 border-sky-400 ml-[10px] rounded-md hover:bg-slate-400"
                >
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
