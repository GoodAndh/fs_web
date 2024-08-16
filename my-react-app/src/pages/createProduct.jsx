import { useState, useEffect } from "react";
import ModalLoginErrorMsg from "../components/ModalLoginErrorMsg.jsx";
import Form, { BaseForm } from "../components/Form.jsx";
import Navbar from "../components/Navbar.jsx";
import { useAuth } from "../utils/AuthContext.jsx";
import { usePostJson } from "../utils/customHook.js";

function CreateProduct() {
  const { isAuth } = useAuth();
  const [preview1, setPreview1] = useState(null);
  const [preview2, setPreview2] = useState(null);
  const [preview3, setPreview3] = useState(null);
  const defaultForm = {
    name: "",
    description: "",
    price: 0,
    stock: 0,
  };

  const [form, setForm] = useState(defaultForm);
  const {
    data: formResponse,
    error: formError,
    call: formCall,
    resetResponse: formReset,
  } = usePostJson("product/create");
  const { call: imageCall } = usePostJson("product/image");

  const [fileImage1, setFile1] = useState(null);
  const [fileImage2, setFile2] = useState(null);
  const [fileImage3, setFile3] = useState(null);

  const onChange = (id, value, files) => {
    setForm({ ...form, [id]: value });
    if (files && files[0]) {
      const file = files && files[0];

      const reader = new FileReader();
      reader.onloadend = () => {
        if (id === "image1") {
          setPreview1(reader.result);
          setFile1(file && file);
        }
        if (id === "image2") {
          setPreview2(reader.result);
          setFile2(file && file);
        }
        if (id === "image3") {
          setPreview3(reader.result);
          setFile3(file && file);
        }
      };
      reader.readAsDataURL(files[0]);
    } else {
      !preview1
        ? setPreview1(null)
        : !preview2
        ? setPreview2(null)
        : !preview3
        ? setPreview3(null)
        : null;
    }
  };

  const doSubmit = async (e) => {
    if (fileImage1 && fileImage2 && fileImage3) {
      e.preventDefault();
      const formInput = {
        ...form,
        price: parseFloat(form.price),
        stock: parseFloat(form.stock),
      };
      await formCall({ form: formInput });
    }
  };

  async function ImageCreate(files, prID) {
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (file) {
        const fileForm = new FormData();
        fileForm.append("file", file);
        await imageCall({
          form: fileForm,
          addUrls: prID,
          params: {
            url: generateFormattedRandomString(),
            captions: generateFormattedRandomString(),
          },
        });
      }
    }
  }

  useEffect(() => {
    if (formResponse?.data.id) {
      const files = [fileImage1, fileImage2, fileImage3];
      ImageCreate(files, `/${formResponse?.data?.id}`);
      formReset();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [formResponse]);

  function generateFormattedRandomString() {
    const now = new Date();
    const day = String(now.getDate()).padStart(2, "0");
    const month = String(now.getMonth() + 1).padStart(2, "0");
    const year = now.getFullYear();
    const hours = String(now.getHours()).padStart(2, "0");
    const minutes = String(now.getMinutes()).padStart(2, "0");
    const seconds = String(now.getSeconds()).padStart(2, "0");

    const timestamp = `${day}-${month}-${year}_${hours}:${minutes}:${seconds}`;

    const randomString = Math.random().toString(36).substring(2, 16);

    return `${timestamp}_${randomString}`;
  }

  return (
    <>
      {!isAuth ? (
        <ModalLoginErrorMsg />
      ) : (
        <>
          <Navbar />
          <div className="flex max-w-[1400px] m-auto mt-[100px] overflow-hidden">
            <div className="w-full mb-10 mx-auto rounded-2xl  overflow-hidden border-[1px] border-slate-400 ">
              <div className="mb-10">
                {/* input image and preview */}
                <BaseForm
                  id="image1"
                  type="file"
                  title=""
                  inputChange={onChange}
                  validate={""}
                  customClassesForm={`hidden`}
                  newLabel={
                    <div className="ml-[80px]">
                      <div className=" border-[1px] border-slate-400 rounded-md w-[480px] h-[300px]">
                        <img
                          src={preview1}
                          alt="you dont have any images to show, try edit"
                          className=""
                        />
                      </div>
                      <label
                        htmlFor="image1"
                        className=" border-[1px] border-slate-400 rounded-sm mt-2"
                      >
                        {/* preview1 */}
                        Tambah gambar
                      </label>
                    </div>
                  }
                />
                <BaseForm
                  id="image2"
                  type="file"
                  title=""
                  inputChange={onChange}
                  validate={""}
                  customClassesForm={`hidden`}
                  newLabel={
                    <div className="ml-[80px]">
                      <div className="bg-center border-[1px] border-slate-400 rounded-md w-[480px] h-[300px]">
                        <img
                          src={preview2}
                          alt="you dont have any images to show, try edit"
                          className=""
                        />
                      </div>
                      <label
                        htmlFor="image2"
                        className=" border-[1px] border-slate-400 rounded-sm mt-2"
                      >
                        {/* preview1 */}
                        Tambah gambar
                      </label>
                    </div>
                  }
                />
                <BaseForm
                  id="image3"
                  type="file"
                  title=""
                  inputChange={onChange}
                  validate={""}
                  customClassesForm={`hidden`}
                  newLabel={
                    <div className="ml-[80px]">
                      <div className=" border-[1px] border-slate-400 rounded-md w-[480px] h-[300px]">
                        <img
                          src={preview3}
                          alt="you dont have any images to show, try edit"
                          className=""
                        />
                      </div>
                      <label
                        htmlFor={`image3`}
                        className={`border-[1px] border-slate-400 rounded-sm mt-2 `}
                      >
                        {/* preview1 */}
                        Tambah gambar
                      </label>
                    </div>
                  }
                />
              </div>
            </div>
            <div className="bg-greesn-400 w-[50px]">
              {/* batas jangan di isi apapun */}
            </div>
            <div className="w-full h-auto mb-10  mx-auto rounded-2xl border-[1px] border-slate-400  ">
              {/* form produk */}
              <div className="ml-[10px] mt-[30px] w-[500px]">
                <Form
                  id="name"
                  type="text"
                  title="Nama Produk"
                  inputChange={onChange}
                  validate={formError?.data?.Name || ""}
                />
                <br />
                <BaseForm
                  id="description"
                  type="textarea"
                  title="Deskripsi Produk"
                  inputChange={onChange}
                  validate={formError?.data?.Description || ""}
                  customClassesForm={`w-[400px] h-[120px]  resize-none mt-2 px-3 py-2 shadow-md rounded-md cursor-pointer`}
                />
                <br />
                <div className="flex mt-[35px]">
                  <div className="mr-[100px]">
                    <BaseForm
                      id="price"
                      type="number"
                      title="Price"
                      inputChange={onChange}
                      validate={formError?.data?.Price || ""}
                      customClassesForm={`ml-2 text-center w-[100px] resize-none mt-2 px-3 py-2 shadow-md rounded-md cursor-pointer`}
                    />
                  </div>
                  <div className="">
                    <BaseForm
                      id="stock"
                      type="number"
                      title="Quantity"
                      inputChange={onChange}
                      validate={formError?.data?.Stock || ""}
                      customClassesForm={`ml-2 text-center w-[100px] resize-none mt-2 px-3 py-2 shadow-md rounded-md cursor-pointer`}
                    />
                  </div>
                </div>
              </div>
              <div className=" items-center justify-center mt-[50px]">
                <p className="font-bold ml-[100px] w-[450px] underline">
                Pastikan semua form terisi dan gambar minimal 3 gambar
                </p>
                <button
                  onClick={doSubmit}
                  className={`${
                    !fileImage1 || !fileImage2 || !fileImage3 || !form.Name
                      ? "cursor-not-allowed"
                      : ""
                  } mt-[50px] w-3/4 mx-[50px]   p-2 bg-transparent hover:bg-green-500 text-slate-700 font-semibold hover:text-white py-2 px-4 border border-slate-500 hover:border-transparent rounded-xl`}
                >
                  Submit
                </button>{" "}
              </div>
              {/*  */}
            </div>
          </div>
        </>
      )}
    </>
  );
}

export default CreateProduct;
