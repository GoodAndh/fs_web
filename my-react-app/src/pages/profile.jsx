/* eslint-disable react/jsx-key */
/* eslint-disable no-prototype-builtins */
import { useEffect, useState } from "react";
import { useAuth } from "../utils/AuthContext.jsx";
import { useGetJson } from "../utils/customHook.js";
import Navbar from "../components/Navbar.jsx";
import Form from "../components/Form.jsx";
import ModalLoginErrorMsg from "../components/ModalLoginErrorMsg.jsx";
import ImageSlider from "../components/ProductImageSlider.jsx";

import { Link } from "react-router-dom";

import { usePostJson } from "../utils/customHook.js";

import FiEdit from "../assets/edit.svg";

import Images from "../components/Images.jsx";
{
  /* <Images
url="user/serveprofile"
params={{ url: "urlUpdateCoba.png " }}
/> */
}
function Profile() {
  const { data: url, error: error } = useGetJson("user/profile");
  const { data: data } = useGetJson("user/getuser");
  const { data: myProduct } = useGetJson("product/ownproduct");

  const { isAuth } = useAuth();
  const [isOpen, setOpen] = useState(false);

  const [file, setFile] = useState(null);
  const [previewUrl, setPreview] = useState(null);
  const formInitial = {
    captions: "",
    url: "",
  };
  const [form, setForm] = useState(formInitial);
  const initialResponse = {
    fileMessage: "",
    Url: "",
    Captions: "",
  };
  const [response, setResponse] = useState(initialResponse);

  const {
    data: createResponse,
    error: errorCreate,
    call: callCreate,
    resetResponse: resetCreate,
  } = usePostJson("user/profile/create");
  const {
    data: updateResponse,
    error: errorUpdate,
    call: callUpdate,
    resetResponse: resetUpdate,
  } = usePostJson("user/profile/update");

 

  const clickBro = () => {
    setOpen((is) => !is);
  };

  const onChange = (id, value, files) => {
    setForm({ ...form, [id]: value });
    files && files[0] && setFile(files[0]);

    const reader = new FileReader();
    reader.onloadend = () => {
      setPreview(reader.result);
    };

    if (files && files[0]) {
      reader.readAsDataURL(files[0]);
    } else {
      setPreview(null);
    }
  };

  const doSubmit = async (e) => {
    e.preventDefault;

    const fileForm = new FormData();
    fileForm.append("file", file);

    await callCreate({
      form: fileForm && fileForm,
      params: { url: form.url, captions: form.captions },
      fallbackFunc: await fbFunc({
        form: fileForm && fileForm,
        params: { url: form.url, captions: form.captions },
      }),
    });
  };

  const fbFunc = async ({ form, params }) => {
    await callUpdate({ form: form, params: params });
  };

  useEffect(() => {
    if (updateResponse) {
      setResponse(initialResponse);
      resetUpdate();
    }
    if (errorUpdate) {
      const newRes = { ...response };
      if (errorUpdate.data) {
        Object.keys(errorUpdate.data).forEach((keys) => {
          response.hasOwnProperty(keys) &&
            (newRes[keys] = errorUpdate.data[keys]);
        });
        Object.keys(response).forEach((keys) => {
          !errorUpdate.data.hasOwnProperty(keys) && (newRes[keys] = "");
        });
      }
      setResponse(newRes);
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [updateResponse, errorUpdate]);

  useEffect(() => {
    if (createResponse) {
      setResponse(initialResponse);
      resetCreate();
    }

    if (errorCreate) {
      const newRes = { ...response };
      if (errorCreate.data) {
        Object.keys(errorCreate.data).forEach((keys) => {
          response.hasOwnProperty(keys) &&
            (newRes[keys] = errorCreate.data[keys]);
        });
        Object.keys(response).forEach((keys) => {
          !errorCreate.data.hasOwnProperty(keys) && (newRes[keys] = "");
        });
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [createResponse, errorCreate]);

  return (
    <>
      {!isAuth ? (
        <ModalLoginErrorMsg />
      ) : (
        <>
          <Navbar />

          <div className="max-w-[1400px] h-[1500px] w-full m-auto py-16 px-4  rounded-sm relative ">
            <div className="w-[700px] h-[390px] mx-auto rounded-2xl overflow-hidden  shadow-lg ">
              {url && !previewUrl && (
                <Images
                  url="user/serveprofile"
                  params={{ url: `${url && url.data.url && url.data.url}` }}
                  captions={`${
                    url && url.data.captions
                      ? url.data.captions
                      : error && error.message
                  }`}
                />
              )}
              {!url ||
                (previewUrl && (
                  <img
                    src={previewUrl}
                    alt="you dont have any images to show try edit"
                    className="ml-[10px]"
                  />
                ))}
            </div>

            <div
              className="ml-[990px] w-[45px] h-[50px] rounded-lg m-2 hover:cursor-pointer bg-slate-300"
              onClick={clickBro}
            >
              <img src={FiEdit} alt="not found images" className="ml-[10px]" />
              <span className="m-2 font-semibold">Edit</span>
            </div>
            <div
              className={`${
                isOpen ? "block" : "hidden"
              } w-[700px] h-[350x] p-4 mx-auto rounded-md bg-slate-200`}
            >
              <Form
                id="img"
                type="file"
                title="Insert Image"
                inputChange={onChange}
                validate={response.fileMessage}
              />
              <Form
                id="url"
                type="text"
                title="Image Name Initial(Url)"
                inputChange={onChange}
                validate={response.Url}
              />
              <Form
                id="captions"
                type="text"
                title="Captions Name Initial(Captions)"
                inputChange={onChange}
                validate={response.Captions}
              />
              <button
                onClick={doSubmit}
                className="my-5 w-full p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl"
              >
                Submit
              </button>
            </div>

            <div className="m-6 w-[700px] mx-auto p-5  rounded-xl shadow-lg">
              <div className="p-2  flex flex-col">
                <div className="flex items-center">
                  <p className="p-2 ">Username:</p>
                  {data ? (
                    <p className="font-bold">{data.data.username}</p>
                  ) : (
                    <p className="font-bold">Halo</p>
                  )}
                </div>
                <div className="flex items-center">
                  <p className="p-2 ">Email:</p>
                  {data ? (
                    <p className="font-bold">{data.data.email}</p>
                  ) : (
                    <p className="font-bold">Halo</p>
                  )}
                </div>
              </div>
            </div>
            {myProduct &&
              myProduct.data &&
              myProduct.data.map((item) => {
                return (
                  <div
                    key={item.id}
                    className="w-[700px] h-auto mx-auto m-6 rounded-xl shadow-lg"
                  >
                    <div className="ml-[655px] w-[45px] h-[50px] rounded-lg m-2 hover:cursor-pointer bg-slate-300">
                      <img
                        src={FiEdit}
                        alt="not found images"
                        className="ml-[10px]"
                      />
                      <span className="m-2 font-semibold">Edit</span>
                    </div>
                    <div className="flex items-center justify-center container w-[200px] ml-4 my-4">
                      <div className="flex gap-8">
                        <div className="rounded-xl shadow-lg bg-slate-300">
                          <div className="m-2 p-1 ">
                            <ImageSlider itemID={item.id} />
                          </div>
                          <br />
                          <h2 className="text-2xl md:text-3xl m-1">
                            {item.name}
                          </h2>
                          <p className="text-md font-medium m-1">
                            {" "}
                            Rp.{item.price}
                          </p>
                          <div className="my-4">
                            <Link
                            to={`/p/${item.name}`}
                            className="my-5 w-full p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl"
                            >
                              Selengkapnya
                            </Link>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
          </div>
        </>
      )}
    </>
  );
}

export default Profile;
